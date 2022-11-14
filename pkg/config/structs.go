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

package config

import (
	"fmt"

	"github.com/creasty/defaults"
)

type Operations struct {
	Add  string `default:"apk add {args}" json:"add"` // add "{-y}" to auto install on creation on debian
	Del  string `default:"apk del {args}" json:"del"`
	Upd  string `default:"apk update {args}" json:"update"`
	Dup  string `default:"apk upgrade {args}" json:"upgrade"`
	Srch string `default:"apk search {args}" json:"search"`
	Clean string `default:"rm -rf /var/cache/apk/*" json:"clean"`
}

type Installer struct {
	Operations  Operations        `json:"operations"`
	ArgModifier map[string]string `default:"{}" json:"args-modifier"`
}

type Image struct {
	URI        string            `default:"alpine:latest" json:"uri"`
	OnCreation []string          `default:"[\"apk update\"]" json:"on-creation"`
	OnFinish   []string          `default:"[]" json:"on-finish"`
	Installer  Installer         `json:"pkg-manager"`
	EnvVars    map[string]string `default:"{}" json:"env-vars"`
}

type Binds struct {
	XOrg bool     `default:"true" json:"xorg"`
	Dev  bool     `default:"true" json:"/dev"`
	Vars []string `default:"[]" json:"env-vars"`
}

type Experiments struct {
	Socket bool `default:"false" json:"sockets"`
}

type Container struct {
	Name        string      `json:"name"`
	Args        []string    `default:"[\"--net=host\"]" json:"arguments"`
	WorkDir     string      `default:"/code" json:"work-dir"`
	Shell       string      `default:"/bin/sh" json:"shell"`
	RootUser    bool        `json:"root-user"`
	Privileged  bool        `default:"true" json:"privileged"`
	Binds       Binds       `json:"binds"`
	Ports       []string    `default:"[]" json:"ports"`
	Mounts      []string    `default:"[]" json:"mounts"`
	Experiments Experiments `json:"experiments"`
	// This is a map of tag:container_path
	// For example, for nix, you can use
	// the nix tag (which will create a global folder called nix)
	// and then use it in the container as /nix
	SharedFolders map[string]interface{} `default:"{}" json:"shared-folders"`
}

type Podman struct {
	Path      string    `default:"podman" json:"path"`
	Rootless  bool      `default:"true" json:"rootless"`
	BuildOnly bool      `json:"create-deletion"`
	Container Container `json:"container"`
}

type UserPkgs struct {
	Packages    []string `default:"[]" json:"packages"`
	DevPackages []string `default:"[]" json:"devpackages"`
}

type Struct struct {
	Image  Image  `json:"image"`
	Podman Podman `json:"podman"`

	// Using interface so we can support string and []string
	Commands map[string]interface{} `default:"{}" json:"commands"`

	Packages    []string `default:"[]" json:"packages"`
	DevPackages []string `default:"[]" json:"devpackages"`
	UserPkgs    UserPkgs `json:"userpkgs"`
}

func SetName(cfg *Struct) {
	if cfg.Podman.Container.Name == "" {
		cfg.Podman.Container.Name = fmt.Sprintf("develbox-%s", GetDirNmHash()[:32])
	}
}

func SetDefaults(cfg *Struct) {
	defaults.Set(cfg)
	SetName(cfg)
}
