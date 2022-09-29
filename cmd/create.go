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

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	createCfg    bool
	forceReplace bool
	downloadUrl  string
	Create       = &cobra.Command{
		Use:        "create [CONFIG_NAME]",
		SuggestFor: []string{"config"},
		Short:      "Creates a new container/config for this project",
		Args:       cobra.MaximumNArgs(1),
		Example:    "develbox create -c alpine/latest",
		Run: func(cmd *cobra.Command, args []string) {
			configExists := config.ConfigExists()

			if createCfg {
				if configExists && !forceReplace {
					glg.Errorf("Config file already exists!")
					os.Exit(1)
				}

				cfg := config.Struct{}

				if len(args) == 1 {
					cfg = downloadConfig(args[0])
				}

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
				config.SetDefaults(&cfg)

				err = config.Write(&cfg)
				if err != nil {
					glg.Error(err)
				}
				fmt.Println("Config file created!")
				return
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
	Create.Flags().StringVarP(&downloadUrl, "source", "s", "https://raw.githubusercontent.com/kadmuffin/develbox/main/configs", "A base path from where to get the configs.")
}

func downloadConfig(argum string) config.Struct {
	resp, err := http.Get(fmt.Sprintf("%s/%s.json", downloadUrl, argum))

	if err != nil {
		glg.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		glg.Errorf("Response from source returned bad status: %s", resp.StatusCode)
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Fatalf("Something went wrong while downloading the config file. %s", err)
	}

	cfg, err := config.ReadBytes(data)
	if err != nil {
		glg.Fatalf("Failed to parse the JSON data. %s", err)
	}
	return cfg
}
