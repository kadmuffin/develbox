// Copyright 2022 Kevin Ledesma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	// Start is the cobra command for the start command
	Start = &cobra.Command{
		Use:     "start",
		Aliases: []string{"up"},
		Short:   "Starts the container",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Read()
			if err != nil {
				glg.Fatal(err)
			}
			pman := podman.New(cfg.Podman.Path)

			if !pman.Exists(cfg.Container.Name) {
				glg.Fatal("Container does not exist")
			}

			err = StartContainer(cfg.Container.Name, pman, podman.Attach{})
			if err != nil {
				glg.Fatal(err)
			}
			fmt.Println("Container started.")
		},
	}
)

func StartContainer(name string, pman podman.Podman, attach podman.Attach) error {
	containers, err := SearchActiveContainer(pman)
	if err != nil {
		glg.Warn(err)
	}

	// Remove this container from the list
	for i, container := range containers {
		if container.Name == name {
			containers = append(containers[:i], containers[i+1:]...)
		}
	}

	if len(containers) > 0 {
		prompt := promptui.Select{
			Label: "Found other active containers, do you want to stop them?",
			Items: []string{"Stop them", "Choose what to stop", "Don't stop them"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			return err
		}

		switch result {
		case "Stop them":
			for _, container := range containers {
				err := pman.Stop([]string{container.Name}, podman.Attach{})
				if err != nil {
					glg.Warn(err)
				}
			}
		case "Choose what to stop":
			ChooseWhatToStop(pman, containers)
		default:
			fmt.Println("Not stopping any containers")
		}
	}
	return pman.Start([]string{name}, attach)
}

// ChooseWhatToStop is a prompt that allows the user to choose what to stop
// It's a recursive function that doesn't exit until the user selects "Done"
func ChooseWhatToStop(pman podman.Podman, containers []contInfo) error {
	contString := make([]string, len(containers))
	for i, container := range containers {
		contString[i] = container.String()
	}
	prompt := promptui.Select{
		Label: "Select a container to stop",
		Items: contString,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}

	glg.Debugf("Choosed container: %s", result)
	for i, container := range containers {
		glg.Debugf("Current container on list (string): %s", container.String())
		if container.String() == result {
			fmt.Println("Stopping container ", container.Name)
			err := pman.Stop([]string{container.Name}, podman.Attach{
				Stdout: true,
				Stderr: true,
			})
			if err != nil {
				glg.Warn(err)
			}

			// Remove the container from the list
			containers = append(containers[:i], containers[i+1:]...)
		}
	}

	var items = []string{}
	if len(containers) > 0 {
		items = append(items, "Choose another container to stop")
	}
	items = append(items, "Done")
	prompt = promptui.Select{
		Label: "Select an option",
		Items: items,
	}

	_, result, err = prompt.Run()

	if err != nil {
		return err
	}

	switch result {
	case "Stop another container":
		return ChooseWhatToStop(pman, containers)
	default:
		return nil
	}
}

var keyval = regexp.MustCompile(`((.*)=(.*))`)

// contInfo is a struct to store the information of a container
type contInfo struct {
	Name        string
	ID          string
	ProjectPath string
	DBoxVersion string
}

func (e *contInfo) String() string {
	return fmt.Sprintf("%s (%s) Path: %s", e.Name, e.ID, e.ProjectPath)
}

// SearchActiveContainer searches for all active containers and returns id, name, and project path of containers that match
// Mainly, it searches containers with the label develbox_container=1
// For the project path it gets the label develbox_project_path
func SearchActiveContainer(pman podman.Podman) ([]contInfo, error) {
	cmd := pman.RawCommand([]string{"ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Labels}}"}, podman.Attach{
		Stderr: true,
	})
	containers, err := cmd.Output()
	if err != nil {
		return []contInfo{}, err
	}

	// Split the output into lines
	lines := strings.Split(string(containers), "\n")

	// Create a slice to store the containers
	var activeContainers []contInfo

	// Loop through the lines
	for _, line := range lines {
		// Split the line into fields
		fields := strings.Split(line, "\t")

		if container.ContainsString(fields, "develbox") {
			activeContainers = append(activeContainers, contInfo{
				Name:        fields[1],
				ID:          fields[0],
				ProjectPath: parseFromField(fields[2], "develbox_project_path"),
				DBoxVersion: parseFromField(fields[2], "develbox_version"),
			})
		}
	}

	return activeContainers, nil
}

// parseFromField parses a field from a podman ps --format
func parseFromField(field string, key string) string {
	// Remove the curly braces
	field = strings.Trim(field, "{}")
	// Split the field into key value pairs
	pairs := strings.Split(field, ",")
	// Loop through the pairs
	for _, pair := range pairs {
		// Check if the key matches
		if strings.Contains(pair, key) {
			// Split the pair into key and value
			split := keyval.FindStringSubmatch(pair)
			// Return the value
			return split[3]
		}
	}
	return ""
}
