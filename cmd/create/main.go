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

// Package create contains the create command
package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/kadmuffin/develbox/cmd/version"
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kpango/glg"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

var (
	createCfg      bool
	forceReplace   bool
	downloadURL    string
	containerName  string
	containerMount string
	containerPort  string
	versionTag     string

	// Create is the main command for creating a container
	Create = &cobra.Command{
		Use:        "create",
		SuggestFor: []string{"config", "init"},
		Short:      "Creates a new container/config for this project",
		Args:       cobra.MaximumNArgs(1),
		Example:    "develbox create -c alpine/latest",
		Run: func(cmd *cobra.Command, args []string) {
			configExists := config.Exists()

			if createCfg || !configExists {
				if configExists && !forceReplace {
					glg.Errorf("Config file already exists!\nUse -f to force the creation of a new config file.")
					os.Exit(1)
				}

				cfg := config.Structure{}

				switch isURL(downloadURL) {
				case true:
					switch len(args) {
					case 0:
						targetURL := strings.ReplaceAll(downloadURL, "$$version$$", "v"+versionTag)
						fmt.Println(targetURL)
						cfg = promptConfig(targetURL)
					default:
						cfg = downloadConfig(args[0], strings.ReplaceAll(downloadURL, "$$tag$$", "main"))
					}
				case false:
					var err error
					var v1Cfg bool
					cfg, v1Cfg, err = config.ReadFile(downloadURL)
					if err != nil {
						glg.Fatalf("Couldn't read config file: %s", err)
					}

					if v1Cfg {
						glg.Warn("Config file is from an older version of develbox. Develbox will write an updated version to .develbox/config.json")
					}
				}

				config.SetDefaults(&cfg)

				if containerName == "" {
					containerName = promptName(&cfg)
				}
				cfg.Container.Name = containerName

				if containerMount == "none" {
					containerMount = promptVolumes(&cfg)
				}

				switch containerMount {
				case "":
					cfg.Container.Mounts = []string{}
				default:
					cfg.Container.Mounts = strings.Split(containerMount, ",")
				}

				if containerPort == "none" {
					containerPort = promptPorts(&cfg)
				}

				switch containerPort {
				case "":
					cfg.Podman.Args = append(cfg.Podman.Args, "--net=host")
				default:
					cfg.Container.Ports = strings.Split(containerPort, ",")
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
			container.PkgVersion = cmd.Root().Version
			container.Create(cfg, forceReplace)
		},
	}
)

func init() {
	Create.Flags().BoolVarP(&createCfg, "config", "c", false, "Use to create a new config file")
	Create.Flags().BoolVarP(&forceReplace, "force", "f", false, "Use to force the creation of a container/config")
	Create.Flags().StringVarP(&downloadURL, "source", "s", "https://raw.githubusercontent.com/kadmuffin/develbox/$$version$$/configs", "A base path from where to get the configs.")
	Create.Flags().StringVarP(&containerName, "name", "n", "", "The name of the container to create.")
	Create.Flags().StringVarP(&containerMount, "mount", "m", "none", "The volume to mount in the container.")
	Create.Flags().StringVarP(&containerPort, "port", "p", "none", "The port to expose in the container.")
	Create.Flags().StringVarP(&versionTag, "version", "v", version.Number, "The version tag from where the config will be downloaded.")

}

// setupGitIgnore adds the .develbox directory to the .gitignore file (if the user accepts)
func setupGitIgnore() error {
	gitign, _ := ignore.CompileIgnoreFile(".gitignore")

	matches := (gitign.MatchesPath(".develbox/home/") || gitign.MatchesPath(".develbox/"))

	if !matches && promptGitignore() {
		return writeGitIgnore()
	}
	return nil
}

// writeGitignore writes to the .gitignore file ".develbox/home/"
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

// isURL checks if it's a valid URL or a file path
func isURL(URL string) bool {
	return strings.HasPrefix(URL, "http://") || strings.HasPrefix(URL, "https://")
}
