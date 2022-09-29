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
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

// A struct that is used to install/delete packages.
type operation struct {
	Type        string   `json:"operation"`
	Packages    []string `json:"packages"`
	Flags       []string `json:"flags"`
	AutoInstall bool     `json:"auto-install"`
}

// Creates a new operation struct that is used to request a transaction.
//
// After creating a operation, run yourOperation.Process() to install/delete/etc
//
// Accepted types: ("add", "del", "update", "upgrade", "search")
func NewOperation(opType string, packages []string, flags []string, autoInstall bool) operation {
	return operation{Type: opType, Packages: packages, Flags: flags, AutoInstall: autoInstall}
}

// Processes the transaction and updates the config reference.
//
// Returns an error in case of failure.
func (e *operation) Process(cfg *config.Struct) error {
	pman := podman.New(cfg.Podman.Path)
	cname := cfg.Podman.Container.Name

	switch e.Type {
	case "add":
		if err := e.sendCommand(cname, cfg.Image.Installer.Operations.Add, pman, cfg.Image.Installer.ArgModifier["add"]).Run(); err != nil {
			return glg.Errorf("couldn't install the requested packages: %s", err)
		}

		cfg.Packages = append(cfg.Packages, RemoveDuplicates(&cfg.Packages, &e.Packages)...)

		return nil

	case "del":
		if err := e.sendCommand(cname, cfg.Image.Installer.Operations.Del, pman, cfg.Image.Installer.ArgModifier["del"]).Run(); err != nil {
			return glg.Errorf("couldn't removing the requested packages: %s", err)
		}

		cfg.Packages = RemoveDuplicates(&e.Packages, &cfg.Packages)

		return nil

	case "search":
		if err := e.sendCommand(cname, cfg.Image.Installer.Operations.Srch, pman, cfg.Image.Installer.ArgModifier["search"]).Run(); err != nil {
			return glg.Errorf("couldn't search the requested packages: %s", err)
		}
		return nil

	case "update":
		if err := e.sendCommand(cname, cfg.Image.Installer.Operations.Upd, pman, cfg.Image.Installer.ArgModifier["update"]).Run(); err != nil {
			return glg.Errorf("something failed while running update: %s", err)
		}
		return nil

	case "upgrade":
		if err := e.sendCommand(cname, cfg.Image.Installer.Operations.Upd, pman, cfg.Image.Installer.ArgModifier["upgrade"]).Run(); err != nil {
			return glg.Errorf("something failed while running update: %s", err)
		}
		return nil
	}

	return glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
}

// Processes the transaction and returns a command
//
// Config updates have to be handle separately
func (e *operation) ProcessCmd(cfg *config.Struct) (*exec.Cmd, error) {
	pman := podman.New(cfg.Podman.Path)
	cname := cfg.Podman.Container.Name

	switch e.Type {
	case "add":
		return e.sendCommand(cname, cfg.Image.Installer.Operations.Add, pman, cfg.Image.Installer.ArgModifier["add"]), nil

	case "del":
		return e.sendCommand(cname, cfg.Image.Installer.Operations.Del, pman, cfg.Image.Installer.ArgModifier["del"]), nil

	case "search":
		return e.sendCommand(cname, cfg.Image.Installer.Operations.Srch, pman, cfg.Image.Installer.ArgModifier["search"]), nil

	case "update":
		return e.sendCommand(cname, cfg.Image.Installer.Operations.Upd, pman, cfg.Image.Installer.ArgModifier["update"]), nil

	case "upgrade":
		return e.sendCommand(cname, cfg.Image.Installer.Operations.Upd, pman, cfg.Image.Installer.ArgModifier["upgrade"]), nil
	}

	return &exec.Cmd{}, glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
}

// Runs a podman command with the config's pkgmanager settings.
//
// Returns an error so we can know if something failed or the user
// did a Ctrl+C and stopped the transaction. Either way, packages failed to install.
func (e *operation) sendCommand(cname, base string, pman podman.Podman, argModifier string) *exec.Cmd {
	packages := processPackages(e.Packages, argModifier)
	flags := strings.Join(e.Flags, " ")
	modifBase := strings.Replace(base, "{args}", fmt.Sprintf("%s %s", packages, flags), 1)
	if e.AutoInstall {
		modifBase = strings.Replace(modifBase, "{-y}", "-y", 1)
	} else {
		modifBase = strings.Replace(modifBase, "{-y}", "", 1)
	}
	arguments := []string{cname, modifBase}

	return pman.Exec(arguments, true, true, podman.Attach{Stdin: true, Stdout: true, Stderr: true})
}

// writes a JSON formatted data into a file.
//
// In this case, it's used to write into the pipe.
func (e *operation) Write(path string, perm os.FileMode) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, perm)
}

// Parses a bytes list and returns aan operation
// and an error.
func Read(data []byte) (operation, error) {
	opertn := operation{}
	err := json.Unmarshal(data, &opertn)
	return opertn, err
}
