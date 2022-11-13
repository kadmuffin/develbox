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
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kpango/glg"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

var (
	createCfg       bool
	forceReplace    bool
	downloadUrl     string
	containerName   string
	containerVolume string
	containerPort   string

	Create = &cobra.Command{
		Use:        "create",
		SuggestFor: []string{"config", "init"},
		Short:      "Creates a new container/config for this project",
		Args:       cobra.MaximumNArgs(1),
		Example:    "develbox create -c alpine/latest",
		Run: func(cmd *cobra.Command, args []string) {
			configExists := config.ConfigExists()

			if createCfg || !configExists {
				if configExists && !forceReplace {
					glg.Errorf("Config file already exists!\nUse -f to force the creation of a new config file.")
					os.Exit(1)
				}

				cfg := config.Struct{}

				switch isUrl(downloadUrl) {
				case true:
					switch len(args) {
					case 0:
						targetUrl := strings.ReplaceAll(downloadUrl, "$$tag$$", "v"+cmd.Root().Version)
						fmt.Println(targetUrl)
						cfg = promptConfig(targetUrl)
					default:
						cfg = downloadConfig(args[0], strings.ReplaceAll(downloadUrl, "$$tag$$", "main"))
					}
				case false:
					var err error
					cfg, err = config.ReadPath(downloadUrl)
					if err != nil {
						glg.Fatalf("Couldn't read config file: %s", err)
					}
				}

				checkDocker(&cfg)
				config.SetDefaults(&cfg)

				if containerName == "" {
					containerName = promptName(&cfg)
				}
				cfg.Podman.Container.Name = containerName

				if containerVolume == "" {
					containerVolume = promptVolumes(&cfg)
				}
				cfg.Podman.Container.Mounts = strings.Split(containerVolume, ",")

				if containerPort == "" {
					containerPort = promptPorts(&cfg)
				}

				switch containerPort {
				case "":
					cfg.Podman.Container.Args = append(cfg.Podman.Container.Args, "--net=host")
				default:
					cfg.Podman.Container.Ports = strings.Split(containerPort, ",")
				}

				err := config.Write(&cfg)
				if err != nil {
					glg.Error(err)
				}

				fmt.Println("Config file created!")

				promptEditConfig()

				if config.FileExists(".git") || config.FileExists(".gitignore") {
					err = setupGitIgnore()
					if err != nil {
						glg.Error(err)
					}
				}

				if !promptCont() {
					fmt.Println("Run again the following command create when you're ready to create the container:")
					fmt.Println("develbox create")
					return
				}
			}

			cfg, err := config.Read()
			if err != nil {
				glg.Errorf("Failed to read .develbox/config.json! Try running 'develbox create -c --force' to create a new one.")
				return
			}
			container.Create(cfg, forceReplace)
		},
	}
)

func init() {
	Create.Flags().BoolVarP(&createCfg, "config", "c", false, "Use to create a new config file")
	Create.Flags().BoolVarP(&forceReplace, "force", "f", false, "Use to force the creation of a container/config")
	Create.Flags().StringVarP(&downloadUrl, "source", "s", "https://raw.githubusercontent.com/kadmuffin/develbox/$$tag$$/configs", "A base path from where to get the configs.")
	Create.Flags().StringVarP(&containerName, "name", "n", "", "The name of the container to create.")
	Create.Flags().StringVarP(&containerVolume, "volume", "v", "", "The volume to mount in the container.")
	Create.Flags().StringVarP(&containerPort, "port", "p", "", "The port to expose in the container.")
}

func checkDocker(cfg *config.Struct) {
	err := exec.Command(cfg.Podman.Path, "--version").Run()
	if err != nil {
		err = exec.Command("docker", "--version").Run()
		if err == nil {
			cfg.Podman.Path = "docker"
			glg.Warn("Couldn't find podman! Using docker instead.")
		} else {
			glg.Warn("Couldn't find podman nor docker on PATH!")
		}
	}
}

func setupGitIgnore() error {
	gitign, _ := ignore.CompileIgnoreFile(".gitignore")

	matches := (gitign.MatchesPath(".develbox/home/") || gitign.MatchesPath(".develbox/"))

	if !matches && promptGitignore() {
		return writeGitIgnore()
	}
	return nil
}

// promptGitignore prompts the user if they want to add the .develbox/home folder to the .gitignore file
func writeGitIgnore() error {
	toIgnore := "\n.develbox/home\n"
	if !config.FileExists(".gitignore") {
		os.Create(".gitignore")
		toIgnore = ".develbox/home\n"
	}
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(toIgnore); err != nil {
		return err
	}
	return nil
}

// Checks if it's a valid url or a file path
func isUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
