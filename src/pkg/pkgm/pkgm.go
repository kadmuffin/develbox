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
	"os"
	"strings"

	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/kpango/glg"
)

// A struct that is used to install/delete packages.
type operation struct {
	Type     string   `json:"operation"`
	Packages []string `json:"packages"`
	Flags    []string `json:"flags"`
}

// Creates a new operation struct that is used to request a transaction.
//
// After creating a operation, run yourOperation.Process() to install/delete/etc
//
// Accepted types: ("add", "del", "update", "upgrade", "search")
func NewOperation(opType string, packages []string, flags []string) operation {
	return operation{Type: opType, Packages: packages, Flags: flags}
}

// Processes the transaction and updates the config reference.
//
// Returns an error in case of failure.
func (e *operation) Process(cfg *config.Struct) error {
	pman := podman.New(cfg.Podman.Path)

	switch e.Type {
	case "add":
		if err := e.sendCommand(cfg.Image.Installer.Add, pman); err != nil {
			return glg.Errorf("couldn't install the requested packages: %w", err)
		}

		cfg.Packages = append(cfg.Packages, RemoveDuplicates(&cfg.Packages, &e.Packages)...)

		return nil

	case "del":
		if err := e.sendCommand(cfg.Image.Installer.Del, pman); err != nil {
			return glg.Errorf("couldn't removing the requested packages: %w", err)
		}

		cfg.Packages = RemoveDuplicates(&e.Packages, &cfg.Packages)

		return nil

	case "search":
		if err := e.sendCommand(cfg.Image.Installer.Srch, pman); err != nil {
			return glg.Errorf("couldn't search the requested packages: %w", err)
		}
		return nil

	case "update":
		if err := e.sendCommand(cfg.Image.Installer.Upd, pman); err != nil {
			return glg.Errorf("something failed while running update: %w", err)
		}
		return nil

	case "upgrade":
		if err := e.sendCommand(cfg.Image.Installer.Upd, pman); err != nil {
			return glg.Errorf("something failed while running update: %w", err)
		}
		return nil
	}

	return glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
}

// Runs a podman command with the config's pkgmanager settings.
//
// Returns an error so we can know if something failed or the user
// did a Ctrl+C and stopped the transaction. Either way, packages failed to install.
func (e *operation) sendCommand(base string, pman podman.Podman) error {
	arguments := []string{strings.Replace(base, "{args}", "", 1)}
	arguments = append(arguments, e.Packages...)
	arguments = append(arguments, e.Flags...)

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
