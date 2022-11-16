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
	"os"
	"strings"

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
			if !pman.Exists(cfg.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Container.Name}, podman.Attach{})
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

// SearchActiveContainer searches for all active containers and returns id, name, and project path of containers that match
// Mainly, it searches containers with the label develbox_container=1
// For the project path it gets the label develbox_project_path
func SearchActiveContainer(pman podman.Podman) ([]string, error) {
	cmd := pman.RawCommand([]string{"ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Labels}}"}, podman.Attach{
		Stderr: true,
	})
	containers, err := cmd.Output()
	if err != nil {
		return []string{"", "", ""}, err
	}

	// Split the output into lines
	lines := strings.Split(string(containers), "\n")

	// Remove the last line because it is empty
	lines = lines[:len(lines)-1]

	// Create a slice to store the containers
	var activeContainers []string

	// Loop through the lines
	for _, line := range lines {
		// Split the line into fields
		fields := strings.Split(line, "\t")

		// Check if the container is active
		if fields[2] == "develbox_container=1" {
			// Add the container to the slice
			activeContainers = append(activeContainers, fields[0])
			activeContainers = append(activeContainers, fields[1])
			activeContainers = append(activeContainers, fields[2])
		}
	}

	return activeContainers, nil
}
