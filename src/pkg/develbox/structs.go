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

package develbox

type Installer struct {
	Add  string `default:"apk add {args}" json:"add"` // add "{-y}" to auto install on creation on debian
	Del  string `default:"apk del {args}" json:"del"`
	Upd  string `default:"apk update" json:"update"`
	Dup  string `default:"apk upgrade" json:"upgrade"`
	Srch string `default:"apk search" json:"search"`
}

type Image struct {
	URI        string    `default:"alpine:latest" json:"uri"`
	OnCreation []string  `default:"[\"apk update\"]" json:"on-creation"`
	OnFinish   []string  `default:"[]" json:"on-finish"`
	Installer  Installer `json:"pkg-manager"`
}

type Binds struct {
	Wayland    bool `default:"true" json:"wayland"`
	XOrg       bool `default:"true" json:"xorg"`
	Pulseaudio bool `default:"true" json:"pulseaudio"`
	Pipewire   bool `json:"pipewire"`
	DRI        bool `default:"true" json:"dri"`
	Camera     bool `default:"true" json:"camera"`
}

type Container struct {
	Name     string `json:"name"`
	Args     string `default:"--net=host" json:"arguments"`
	WorkDir  string `default:"/code" json:"work-dir"`
	Shell    string `default:"/bin/sh" json:"shell"`
	RootUser bool   `json:"root-user"`
	Binds    Binds
	Ports    []string `default:"[]" json:"ports"`
	Mounts   []string `default:"[]" json:"mounts"`
}

type Podman struct {
	Path      string    `default:"podman" json:"path"`
	Rootless  bool      `default:"true" json:"rootless"`
	BuildOnly bool      `json:"create-deletion"`
	Container Container `json:"container"`
}

type DevSetings struct {
	Image    Image               `json:"image"`
	Podman   Podman              `json:"podman"`
	Commands map[string][]string `default:"{}" json:"commands"`
	Packages []string            `default:"[]" json:"packages"`
}
