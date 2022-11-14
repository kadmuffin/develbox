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

// Package cmd contains the some commands for the program
package cmd

import (
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	root bool

	// Enter is the cobra command for the enter command
	Enter = &cobra.Command{
		Use:     "enter",
		Aliases: []string{"shell"},
		Short:   "Launches a shell inside the container",
		Long: `Launches a the shell defined in the config inside the container.
		
		To install packages inside the container use the develbox`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			cfg, err := config.Read()
			if err != nil {
				glg.Failf("Can't read config: %s", err)
			}
			pman := podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})
			if socketExperiment && !root {
				go createSocket(&cfg)
			}
			defer os.Remove(".develbox/home/.develbox.sock")
			err = container.InstallAndEnter(cfg, root)
			return err
		},
	}
)

func init() {
	Enter.Flags().BoolVarP(&root, "root", "r", false, "Use to start a root shell in the container")
}
