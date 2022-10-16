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

package create

import (
	"fmt"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
	"github.com/manifoldco/promptui"
)

var avaliableCfgs = map[string][]string{
	"debian": {
		"testing",
		"bullseye",
	},
	"ubuntu": {
		"jammy",
	},
	"fedora": {
		"rawhide",
		"37",
	},
	"arch": {
		"rolling",
	},
	"alpine": {
		"edge",
	},
	"opensuse": {
		"tumbleweed",
	},
	"nix": {
		"unstable",
	},
}

// Create a prompt that asks to choose a distro and version
func promptConfig() config.Struct {
	prompt := promptui.Select{
		Label: "Choose a distro",
		Items: getKeys(avaliableCfgs),
	}
	_, distro, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	prompt = promptui.Select{
		Label: "Choose a version",
		Items: avaliableCfgs[distro],
	}
	_, version, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	return downloadConfig(fmt.Sprintf("%s/%s", distro, version))
}

// Prompts the user for ports to publish
//
// A config reference is required as this function
// will set the network to host if the user chooses
func promptPorts(cfg *config.Struct) {
	prompt := promptui.Prompt{
		Label: "Publish ports (separate with commas, empty for host network)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	if result == "" {
		cfg.Podman.Container.Args = append(cfg.Podman.Container.Args, "--net=host")
		return
	}

	cfg.Podman.Container.Ports = append(cfg.Podman.Container.Ports, strings.Split(result, ",")...)
}

// Prompts the user for volumes to mount
func promptVolumes(cfg *config.Struct) {
	prompt := promptui.Prompt{
		Label: "Mount volumes (separate with commas)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	if result != "" {
		cfg.Podman.Container.Mounts = strings.Split(result, ",")
	}
}

// Prompts the user for the container name
func promptName(cfg *config.Struct) {
	prompt := promptui.Prompt{
		Label: "Container name (empty for default)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	if result != "" {
		cfg.Podman.Container.Name = result
	}
}
