// Copyright 2022 Kevin Ledesma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pkgm manages package installations and removals.
package pkgm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

// Operation is a struct that is used to install/delete packages.
type Operation struct {
	Type          string   `json:"operation"`
	Packages      []string `json:"packages"`
	Flags         []string `json:"flags"`
	AutoInstall   bool     `json:"auto-install"`
	DevInstall    bool     `json:"dev-install"`
	UserOperation bool     `json:"run-as-user"`
}

// NewOperation creates a new operation struct that is used to request a transaction. Accepted types: ("add", "del", "update", "upgrade", "search", "clean").
func NewOperation(opType string, packages []string, flags []string, autoInstall bool) Operation {
	return Operation{Type: opType, Packages: packages, Flags: flags, AutoInstall: autoInstall, UserOperation: false}
}

// UpdateConfig updates the config file with the new packages
func (e *Operation) UpdateConfig(cfg *config.Structure) {
	pkgsP := &cfg.Packages
	devPkgsP := &cfg.DevPackages
	if e.UserOperation {
		pkgsP = &cfg.UserPkgs.Packages
		devPkgsP = &cfg.UserPkgs.DevPackages
	}

	if pkgsP == nil {
		*pkgsP = []string{}
	}
	if devPkgsP == nil {
		*devPkgsP = []string{}
	}

	switch e.Type {
	case "add":
		*pkgsP = RemoveDuplicates(&e.Packages, pkgsP)
		*devPkgsP = RemoveDuplicates(&e.Packages, devPkgsP)

		switch e.DevInstall {
		case true:
			*devPkgsP = append(*devPkgsP, e.Packages...)
		case false:
			*pkgsP = append(*pkgsP, e.Packages...)
		}

	case "del":
		*pkgsP = RemoveDuplicates(&e.Packages, pkgsP)
		*devPkgsP = RemoveDuplicates(&e.Packages, devPkgsP)
	}
}

// Process processes the transaction and updates the config reference. Returns an error in case of failure.
func (e *Operation) Process(cfg *config.Structure) error {
	cmd, err := e.ProcessCmd(cfg, podman.Attach{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		PseudoTTY: true,
	})
	if err != nil {
		return err
	}

	e.UpdateConfig(cfg)

	return cmd.Run()
}

// ProcessCmd processes the transaction and returns a command. Config updates have to be handle separately.
func (e *Operation) ProcessCmd(cfg *config.Structure, attach podman.Attach) (*exec.Cmd, error) {
	var pman podman.Podman
	if !podman.InsideContainer() && os.Getuid() != 0 {
		pman = podman.New(cfg.Podman.Path)
	}
	cname := cfg.Container.Name
	baseCmd, err := e.StringCommand(&cfg.Image.PkgManager)

	if err != nil {
		return nil, err
	}

	glg.Infof("Creating command to install packages: %s", baseCmd)
	return e.sendCommand(
		cname,
		baseCmd,
		pman,
		attach), nil
}

// StringCommand returns the string that will be send to
// the container. For example: "apt install -y vim".
func (e *Operation) StringCommand(cfg *config.PackageManager) (string, error) {
	var baseCmd string

	switch e.Type {
	case "add":
		baseCmd = cfg.Operations.Add
	case "del":
		baseCmd = cfg.Operations.Del
	case "search":
		baseCmd = cfg.Operations.Srch
	case "update":
		baseCmd = cfg.Operations.Upd
	case "upgrade":
		baseCmd = cfg.Operations.Upg
	case "clean":
		baseCmd = cfg.Operations.Clean

	// Throw an error if the operation is not supported
	default:
		return "", glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
	}

	// We want to process the flags and packages, but only if they are not empty
	flags := strings.Join(e.Flags, " ")
	if flags != "" {
		flags += " "
	}
	packages := processPackages(e.Packages, cfg.Modifiers[e.Type])

	// This is where we replace the "{args}" string with the flags and packages
	modifBase := strings.Replace(baseCmd, "{args}", fmt.Sprintf("%s%s", flags, packages), 1)
	regex := regexp.MustCompile(`\{(.*)\}`)

	// Anything that is inside {} will be replaced
	// at this point.
	// If auto install is true, we will remove the brackets
	// else, we just replace them.
	if e.AutoInstall {
		modifBase = regex.ReplaceAllString(modifBase, "$1")
	} else {
		modifBase = regex.ReplaceAllString(modifBase, "")
	}

	return modifBase, nil
}

// sendCommand runs a podman command with the config's pkgmanager settings.
func (e *Operation) sendCommand(cname, base string, pman podman.Podman, attach podman.Attach) *exec.Cmd {

	arguments := []string{cname, base}

	if podman.InsideContainer() && os.Getuid() == 0 {
		if e.UserOperation {
			glg.Warn("Running as root inside a container, but the operation is set to run as user. Ignoring the flag.")
		}

		arguments = strings.Split(base, " ")
		// Because we are inside the container, and we
		// are root, we can just run the command.
		cmd := exec.Command(arguments[0], arguments[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd
	}

	return pman.Exec(arguments, map[string]string{}, true, !e.UserOperation, podman.Attach{Stdin: true, Stdout: true, Stderr: true})
}

// Write writes a JSON formatted data into a file. In this case, it's used to write into the pipe or socket.
func (e *Operation) Write(path string, perm os.FileMode) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, perm)
}

// Read parses a bytes list and returns aan operation and an error.
func Read(data []byte) (Operation, error) {
	opertn := Operation{}
	err := json.Unmarshal(data, &opertn)
	return opertn, err
}

// String returns a string representation of the operation
func (e *Operation) String() string {
	return fmt.Sprintf("Type: %s, Packages: %s, Flags: %s, AutoInstall: %t, DevInstall: %t, UserOperation: %t", e.Type, e.Packages, e.Flags, e.AutoInstall, e.DevInstall, e.UserOperation)
}

// ToJSON converts the operation to a JSON string
func (e *Operation) ToJSON() string {
	data, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(data)
}
