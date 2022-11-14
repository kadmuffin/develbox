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

// Package config contains the configuration for the program
package config

import (
	"fmt"

	"github.com/creasty/defaults"
)

// Operations is a list of package manager commands
type Operations struct {
	// Add is the command to add a package
	Add string `default:"apk add {args}" json:"add"` // add "{-y}" to auto install on creation on debian

	// Del is the command to remove a package
	Del string `default:"apk del {args}" json:"del"`

	// Upd is the command to update a package or the database
	Upd string `default:"apk update {args}" json:"update"`

	// Dup is the command to upgrade a package or all packages
	Dup string `default:"apk upgrade {args}" json:"upgrade"`

	// Srch is the command to search for a package
	Srch string `default:"apk search {args}" json:"search"`

	// Clean is the command to clean the cache
	Clean string `default:"rm -rf /var/cache/apk/*" json:"clean"`
}

type Installer struct {
	// Operations contains the commands for the package manager
	Operations Operations `json:"operations"`

	// ArgModifier lets prefix or suffix packages with a string (see configs/nix/unstable.json)
	ArgModifier map[string]string `default:"{}" json:"args-modifier"`
}

// Image contains the information for the image
type Image struct {
	// URI is the location of the image
	URI string `default:"alpine:latest" json:"uri"`

	// OnCreation is a list of commands to run on creation
	OnCreation []string `default:"[\"apk update\"]" json:"on-creation"`

	// OnFinish is a list of commands to run on finish
	OnFinish []string `default:"[]" json:"on-finish"`

	// Installer contains the configuration for the package manager
	Installer Installer `json:"pkg-manager"`

	// EnvVars is a list of environment variables to set
	EnvVars map[string]string `default:"{}" json:"env-vars"`
}

// Binds is a list of bind mounts
type Binds struct {
	// XOrg decides if the X11 socket should be bind mounted
	XOrg bool `default:"true" json:"xorg"`

	// Dev decides if the /dev directory should be bind mounted
	Dev bool `default:"true" json:"/dev"`

	// Vars is a list of environment variables to copy to the container
	Vars []string `default:"[]" json:"env-vars"`
}

// Experiments contains settings for experimental features
type Experiments struct {
	// Socket decides if the socket server should be started
	Socket bool `default:"false" json:"sockets"`
}

// Container is the struct for the container configuration
type Container struct {
	// Name of the container
	Name string `json:"name"`

	// Args is a list of arguments to pass to the container
	Args []string `default:"[\"--net=host\"]" json:"arguments"`

	// WorkDir is the directory where the code will be mounted
	WorkDir string `default:"/code" json:"work-dir"`

	// Shell is the shell to use in the container
	Shell string `default:"/bin/sh" json:"shell"`

	// RootUser decides if the user inside the container is root or not
	RootUser bool `json:"root-user"`

	// Privileged is a boolean that determines if the container should be run in privileged mode.
	Privileged bool `default:"true" json:"privileged"`

	// Binds contains settings related to the binds (for example, /dev)
	Binds Binds `json:"binds"`

	// Ports is a map of host:container ports
	Ports []string `default:"[]" json:"ports"`

	// Mounts is a list of mounts to be added to the container
	Mounts []string `default:"[]" json:"mounts"`

	// Experiments contains experimental features
	Experiments Experiments `json:"experiments"`

	// SharedFolders is a list of folders that are shared between containers
	SharedFolders map[string]interface{} `default:"{}" json:"shared-folders"`
}

// Podman is the struct for the podman configuration
type Podman struct {
	// Path is the path to the podman executable
	Path string `default:"podman" json:"path"`

	// Rootless tells develbox if podman is running as rootless
	Rootless bool `default:"true" json:"rootless"`

	// BuildOnly tells develbox if it should delete the container after created
	BuildOnly bool `json:"create-deletion"`

	// Container is the container configuration
	Container Container `json:"container"`
}

// UserPkgs is a list of packages to install user-side
type UserPkgs struct {
	// Packages is a list of packages to install
	Packages []string `default:"[]" json:"packages"`

	// DevPackages is a list of packages to install for development
	DevPackages []string `default:"[]" json:"devpackages"`
}

// Struct is the main configuration struct
type Struct struct {
	// Image contains the information for the image
	Image Image `json:"image"`

	// Podman contains the configuration for podman
	Podman Podman `json:"podman"`

	// Using interface so we can support string and []string
	Commands map[string]interface{} `default:"{}" json:"commands"`

	// Packages is a list of packages to install (that are required for your code to work)
	Packages []string `default:"[]" json:"packages"`

	// DevPackages is a list of packages to install for development (like compilers, linters, etc)
	DevPackages []string `default:"[]" json:"devpackages"`

	// UserPkgs is a list of packages to install user-side (Image has to support it)
	UserPkgs UserPkgs `json:"userpkgs"`
}

// SetName sets the name of the container
func SetName(cfg *Struct) {
	if cfg.Podman.Container.Name == "" {
		cfg.Podman.Container.Name = fmt.Sprintf("develbox-%s", GetDirNmHash()[:32])
	}
}

// SetDefaults sets the default values for the configuration
func SetDefaults(cfg *Struct) {
	defaults.Set(cfg)
	SetName(cfg)
}
