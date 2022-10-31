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

// A struct that is used to install/delete packages.
type Operation struct {
	Type          string   `json:"operation"`
	Packages      []string `json:"packages"`
	Flags         []string `json:"flags"`
	AutoInstall   bool     `json:"auto-install"`
	DevInstall    bool     `json:"dev-install"`
	UserOperation bool     `json:"run-as-user"`
}

// Creates a new operation struct that is used to request a transaction.
//
// After creating a operation, run yourOperation.Process() to install/delete/etc
//
// Accepted types: ("add", "del", "update", "upgrade", "search")
func NewOperation(opType string, packages []string, flags []string, autoInstall bool) Operation {
	return Operation{Type: opType, Packages: packages, Flags: flags, AutoInstall: autoInstall, UserOperation: false}
}

// Updates the config file with the new packages
func (e *Operation) UpdateConfig(cfg *config.Struct) {
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

// Processes the transaction and updates the config reference.
//
// Returns an error in case of failure.
// Wrapper around ProcessCmd() that updates
// the config reference.
func (e *Operation) Process(cfg *config.Struct) error {
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

// Processes the transaction and returns a command
//
// Config updates have to be handle separately
func (e *Operation) ProcessCmd(cfg *config.Struct, attach podman.Attach) (*exec.Cmd, error) {
	pman := podman.New(cfg.Podman.Path)
	cname := cfg.Podman.Container.Name
	baseCmd, err := e.StringCommand(&cfg.Image.Installer)

	if err != nil {
		return nil, err
	}

	return e.sendCommand(
		cname,
		baseCmd,
		pman,
		attach), nil
}

// Returns the string that will be send to
// the container.
//
// For example: "apt install -y vim"
func (e *Operation) StringCommand(cfg *config.Installer) (string, error) {
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
		baseCmd = cfg.Operations.Dup

	// Throw an error if the operation is not supported
	default:
		return "", glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
	}

	flags := strings.Join(e.Flags, " ")
	if flags != "" {
		flags += " "
	}
	packages := processPackages(e.Packages, cfg.ArgModifier[e.Type])
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

// Runs a podman command with the config's pkgmanager settings.
//
// Returns an error so we can know if something failed or the user
// did a Ctrl+C and stopped the transaction. Either way, packages failed to install.
func (e *Operation) sendCommand(cname, base string, pman podman.Podman, attach podman.Attach) *exec.Cmd {

	arguments := []string{cname, base}

	return pman.Exec(arguments, map[string]string{}, true, !e.UserOperation, podman.Attach{Stdin: true, Stdout: true, Stderr: true})
}

// writes a JSON formatted data into a file.
//
// In this case, it's used to write into the pipe.
func (e *Operation) Write(path string, perm os.FileMode) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, perm)
}

// Parses a bytes list and returns aan operation
// and an error.
func Read(data []byte) (Operation, error) {
	opertn := Operation{}
	err := json.Unmarshal(data, &opertn)
	return opertn, err
}

// Returns a string representation of the operation
func (e *Operation) String() string {
	return fmt.Sprintf("Type: %s, Packages: %s, Flags: %s, AutoInstall: %t, DevInstall: %t, UserOperation: %t", e.Type, e.Packages, e.Flags, e.AutoInstall, e.DevInstall, e.UserOperation)
}

// Converts the operation to a JSON string
func (e *Operation) ToJSON() string {
	data, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(data)
}
