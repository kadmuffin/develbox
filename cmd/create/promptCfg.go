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
	"os"
	"os/exec"

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
		"tumbleweed_dnf",
	},
	"nix": {
		"unstable",
		"lnl7",
	},
}

// promptConfig create a prompt that asks to choose a distro and version.
func promptConfig(URL string) (config.Structure, error) {
	prompt := promptui.Select{
		Label: "Choose a distro",
		Items: getKeys(avaliableCfgs),
	}
	_, distro, err := prompt.Run()
	if err != nil {
		return config.Structure{}, glg.Errorf("Prompt failed %v\n", err)
	}

	prompt = promptui.Select{
		Label: "Choose a version",
		Items: avaliableCfgs[distro],
	}
	_, version, err := prompt.Run()
	if err != nil {
		return config.Structure{}, glg.Errorf("Prompt failed %v\n", err)
	}

	return downloadConfig(fmt.Sprintf("%s/%s", distro, version), URL)
}

// promptPorts prompts the user for ports to publish.
func promptPorts(cfg *config.Structure) string {
	// A config reference is required as this function
	// will set the network to host if the user chooses
	// to not publish any ports
	prompt := promptui.Prompt{
		Label: "Publish ports (separate with commas, empty for host network)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

// promptVolumes prompts the user for volumes to mount.
func promptVolumes(cfg *config.Structure) string {
	prompt := promptui.Prompt{
		Label: "Mount volumes (separate with commas)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

// promptName prompts the user for the container name.
func promptName(cfg *config.Structure) string {
	prompt := promptui.Prompt{
		Label: "Container name (empty for default)",
	}
	result, err := prompt.Run()
	if err != nil {
		glg.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

// promptGitignore prompts to add .develbox/home to .gitignore. Default answer is yes.
func promptGitignore() bool {
	prompt := promptui.Prompt{
		Label:     "Add .develbox/home to .gitignore",
		IsConfirm: true,
		Default:   "y",
	}
	result, _ := prompt.Run()

	return result == "y"
}

// promptEditConfig asks the user if they want to use $EDITOR to open the config file.
func promptEditConfig() {
	prompt := promptui.Prompt{
		Label:     "Edit config file",
		IsConfirm: true,
		Default:   "y",
	}
	result, _ := prompt.Run()

	if result == "y" || result == "" {
		editor := os.Getenv("EDITOR")
		fmt.Printf("Opening config file in %s\n", editor)

		if editor == "" {
			// Check if at least vi is installed
			_, err := exec.LookPath("vi")
			if err != nil {
				glg.Fatal("No editor set, please set $EDITOR or install vi")
			}
			editor = "vi"
		}

		cmd := exec.Command(editor, ".develbox/config.json")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			glg.Error(err)
		}
	}
}
