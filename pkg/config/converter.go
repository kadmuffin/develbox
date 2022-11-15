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
	"github.com/creasty/defaults"
	v1_config "github.com/kadmuffin/develbox/pkg/config/v1"
)

// ConvertFromV1 converts a v1 config file to a v2 config file
func ConvertFromV1(cfg v1_config.Struct) Struct {
	newCfg := Struct{
		Image: Image{
			URI:        cfg.Image.URI,
			OnCreation: cfg.Image.OnCreation,
			OnFinish:   cfg.Image.OnFinish,
			Variables:  cfg.Image.EnvVars,
			PkgManager: PackageManager{
				Operations: cfg.Image.Installer.Operations,
				Modifiers:  cfg.Image.Installer.ArgModifier,
			},
		},

		Podman: Podman{
			Path:       cfg.Podman.Path,
			Args:       cfg.Podman.Container.Args,
			Privileged: cfg.Podman.Container.Privileged,
			Rootless:   cfg.Podman.Rootless,
			BuildOnly:  cfg.Podman.BuildOnly,
		},

		Container: Container{
			Name:     cfg.Podman.Container.Name,
			WorkDir:  cfg.Podman.Container.WorkDir,
			Shell:    cfg.Podman.Container.Shell,
			RootUser: cfg.Podman.Container.RootUser,
			Binds: Binds{
				XOrg: cfg.Podman.Container.Binds.XOrg,
				Dev:  cfg.Podman.Container.Binds.Dev,
				Vars: cfg.Podman.Container.Binds.Vars,
			},
			Ports:         cfg.Podman.Container.Ports,
			Mounts:        cfg.Podman.Container.Mounts,
			SharedFolders: cfg.Podman.Container.SharedFolders,
		},

		Commands:    cfg.Commands,
		Packages:    cfg.Packages,
		DevPackages: cfg.DevPackages,
		UserPkgs:    cfg.UserPkgs,
		Experiments: cfg.Podman.Container.Experiments,
	}

	defaults.Set(&newCfg)
	SetName(&newCfg)

	return newCfg
}
