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

// Creates a dockerfile based on the config file
package dockerfile

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	command string

	Build = &cobra.Command{
		Use:   "build",
		Short: "Builds a dockerfile based on the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var dckFile []string
			cfg, err := config.Read()
			if err != nil {
				return err
			}

			dckFile = append(dckFile, fmt.Sprintf("FROM %s", cfg.Image.URI))

			// Add precmds before adding any packages
			dckFile = append(dckFile, appendRun(cfg.Image.OnCreation)...)

			// Update base image
			pkgUpdate := pkgm.NewOperation("update", []string{}, []string{}, true)
			update, _ := pkgUpdate.StringCommand(&cfg.Image.Installer)
			dckFile = append(dckFile, fmt.Sprintf("RUN %s", update))

			// Add packages
			pkgInstall := pkgm.NewOperation("add", cfg.Packages, []string{}, true)
			packages, _ := pkgInstall.StringCommand(&cfg.Image.Installer)
			dckFile = append(dckFile, fmt.Sprintf("RUN %s", packages))

			// Add post install commands
			dckFile = append(dckFile, appendRun(cfg.Image.OnFinish)...)

			dckFile = append(dckFile, exposePorts(cfg.Podman.Container.Ports)...)

			for _, mount := range cfg.Podman.Container.Mounts {
				regex := regexp.MustCompile(`([a-zA-Z0-9\.\-\/]+):([a-zA-Z0-9\.\-\/]+):?([a-zA-Z0-9\.\-\/]+)?`)
				match := regex.FindStringSubmatch(mount)
				if len(match) > 0 {
					dckFile = append(dckFile, fmt.Sprintf("COPY %s %s", match[1], match[2]))
				}
			}

			if command != "" {
				// Fail if cfg.Commands map sdoesn't contain the variable command
				if _, ok := cfg.Commands[command]; !ok {
					return glg.Errorf("command %s not found in config file", command)
				}

				dckFile = append(dckFile, fmt.Sprintf("ENTRYPOINT [\"/bin/sh\", \"-c\", \"%s\"]", cfg.Commands[command]))
			} else {
				dckFile = append(dckFile, "ENTRYPOINT [\"/bin/sh\"]")
			}

			if config.FileExists("Dockerfile") {
				glg.Warnf("Dockerfile already exists, overwriting...")
			}
			return writeList("Dockerfile", dckFile)
		},
	}
)

// Opens a file for writing and writes
// a set of lines to it
func writeList(path string, list []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range list {
		_, err := file.WriteString(line + "\n")

		if err != nil {
			return err
		}
	}

	return nil
}

// Parses a port list in the format "host:container"
// and returns a list with the word expose + the container's
// port
func exposePorts(ports []string) []string {
	var exposed []string
	regex := regexp.MustCompile(`[0-9\.]+:+(\d+(\/[a-z]+)?)`)
	for _, port := range ports {
		match := regex.FindStringSubmatch(port)
		if len(match) > 0 {
			exposed = append(exposed, "EXPOSE "+match[1])
		}
	}
	return exposed
}

// Appends the RUN prefix to each element
// in a list.
func appendRun(list []string) []string {
	var newList []string
	for _, line := range list {
		newList = append(newList, "RUN "+line)
	}
	return newList
}

func init() {
	Build.Flags().StringVarP(&command, "command", "c", "", "Command from config file to run on container start")
}
