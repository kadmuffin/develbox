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
	"strings"

	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/kpango/glg"
)

type Operation struct {
	Type     string   `json:"operation"`
	Packages []string `json:"packages"`
	Flags    []string `json:"flags"`
}

func (e *Operation) Process(cfg *config.Struct) error {
	pman := podman.New(cfg.Podman.Path)

	if e.Type == "add" {
		if err := e.sendCommand(cfg.Image.Installer.Add, pman); err != nil {
			return glg.Errorf("couldn't install the requested packages: %w", err)
		}

		cfg.Packages = append(cfg.Packages, RemoveDuplicates(&cfg.Packages, &e.Packages)...)

		return nil
	}

	if e.Type == "del" {
		if err := e.sendCommand(cfg.Image.Installer.Del, pman); err != nil {
			return glg.Errorf("couldn't removing the requested packages: %w", err)
		}

		cfg.Packages = RemoveDuplicates(&e.Packages, &cfg.Packages)

		return nil
	}

	if e.Type == "search" {
		if err := e.sendCommand(cfg.Image.Installer.Srch, pman); err != nil {
			return glg.Errorf("couldn't search the requested packages: %w", err)
		}
		return nil
	}

	if e.Type == "update" {
		if err := e.sendCommand(cfg.Image.Installer.Upd, pman); err != nil {
			return glg.Errorf("something failed while running update: %w", err)
		}
		return nil
	}

	if e.Type == "upgrade" {
		if err := e.sendCommand(cfg.Image.Installer.Upd, pman); err != nil {
			return glg.Errorf("something failed while running update: %w", err)
		}
		return nil
	}

	return glg.Errorf("couldn't find the key '%s' on the list of supported operations", e.Type)
}

func (e *Operation) sendCommand(base string, pman podman.Podman) error {
	arguments := []string{strings.Replace(base, "{args}", "", 1)}
	arguments = append(arguments, e.Packages...)
	arguments = append(arguments, e.Flags...)

	return pman.Exec(arguments, true, true, podman.Attach{Stdin: true, Stdout: true, Stderr: true})
}
