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
	v1_config "github.com/kadmuffin/develbox/pkg/config/v1config"
)

// Operations is a list of pkgm that can be performed on a container
type Operations struct {
	// Add is the base string that the package manager uses to add a package
	Add string `default:"" json:"add"` // add "{-y}" to auto install on creation on debian

	// Del, short for delete, is the base string that the package manager uses to delete a package
	Del string `default:"" json:"del"`

	// Upd, short for update, is the base string that the package manager uses to update the database
	Upd string `default:"" json:"update"`

	// Upg, short for upgrade, is the base string that the package manager uses to upgrade all packages
	Upg string `default:"" json:"upgrade"`

	// Srch, short for search, is the base string that the package manager uses to search for a package
	Srch string `default:"" json:"search"`

	// Clean is the command to clean the cache, necessary if we want to reduce the image size (using the build command)
	Clean string `default:"" json:"clean"`
}

// PackageManager is the configuration for the package manager
type PackageManager struct {
	// Operations contains the commands for the package manager
	Operations Operations `json:"operations"`

	// Modifiers enable prefixs or suffixs on packages names with a string (see configs/nix/unstable.json)
	Modifiers map[string]string `default:"{}" json:"modifiers"`
}

// Image contains the information for the image
type Image struct {
	// URI is the location of the image
	URI string `default:"alpine:latest" json:"uri"`

	// OnCreation is a list of commands to run on creation
	OnCreation []string `default:"[]" json:"on_creation"`

	// OnFinish is a list of commands to run on finish
	OnFinish []string `default:"[]" json:"on_finish"`

	// PkgManager contains the configuration for the package manager
	PkgManager PackageManager `json:"pkgmanager"`

	// Variables is a list of environment variables to set
	Variables map[string]string `default:"{}" json:"variables"`
}

// Binds is a list of bind mounts
type Binds struct {
	// XOrg decides if the X11 socket should be bind mounted
	XOrg bool `default:"true" json:"xorg"`

	// Dev decides if the /dev directory should be bind mounted
	Dev bool `default:"true" json:"dev"`

	// Variables is a list of environment variables to copy to the container
	Variables []string `default:"[]" json:"variables"`
}

// Container is the struct for the container configuration
type Container struct {
	// Name of the container
	Name string `json:"name"`

	// WorkDir is the directory where the code will be mounted
	WorkDir string `default:"/code" json:"workdir"`

	// Shell is the shell to use in the container
	Shell string `default:"/bin/sh" json:"shell"`

	// RootUser decides if the user inside the container is root or not
	RootUser bool `json:"rootuser"`

	// Binds contains settings related to the binds (for example, /dev)
	Binds Binds `json:"binds"`

	// Ports is a map of host:container ports
	Ports []string `default:"[]" json:"ports"`

	// Mounts is a list of mounts to be added to the container
	Mounts []string `default:"[]" json:"mounts"`

	// SharedFolders is a list of folders that are shared between containers
	SharedFolders map[string]interface{} `default:"{}" json:"shared_folders"`
}

// Podman is the struct for the podman configuration
type Podman struct {
	// Path is the path to the podman executable
	Path string `default:"podman" json:"path"`

	// Args is a list of arguments to pass to the podman executable
	Args []string `default:"[]" json:"args"`

	// Rootless tells develbox if podman is running as rootless
	Rootless bool `default:"true" json:"rootless"`

	// AutoDelete tells develbox if it should delete the container after created
	AutoDelete bool `default:"false" json:"onlybuild"`

	// AutoCommit tells develbox if it should commit the container after created
	AutoCommit bool `default:"false" json:"onlycommit"`

	// Privileged is a boolean that determines if the container should be run in privileged mode.
	Privileged bool `default:"true" json:"privileged"`
}

// Structure is the main configuration struct
type Structure struct {
	// Image contains the information for the image
	Image Image `json:"image"`

	// Podman contains the configuration for podman
	Podman Podman `json:"podman"`

	// Container contains the configuration for the container
	Container Container `json:"container"`

	// Using interface so we can support string and []string
	Commands map[string]interface{} `default:"{}" json:"commands"`

	// Packages is a list of packages to install (that are required for your code to work)
	Packages []string `default:"[]" json:"packages"`

	// DevPackages is a list of packages to install for development (like compilers, linters, etc)
	DevPackages []string `default:"[]" json:"devpackages"`

	// UserPkgs is a list of packages to install user-side (Image has to support it)
	UserPkgs v1_config.UserPkgs `json:"userpkgs"`

	// Experiments is a list of experimental features to enable
	Experiments v1_config.Experiments `json:"experiments"`
}

// SetName sets the name of the container
func SetName(cfg *Structure) {
	if cfg.Container.Name == "" {
		cfg.Container.Name = fmt.Sprintf("develbox-%s", v1_config.GetDirNmHash()[:32])
	}
}

// SetDefaults sets the default values for the configuration
func SetDefaults(cfg *Structure) {
	defaults.Set(cfg)
	SetName(cfg)
}
