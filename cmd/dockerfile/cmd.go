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
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kpango/glg"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
)

var (
	command      string
	includeFiles bool
	devBuild     bool

	Build = &cobra.Command{
		Use:   "build",
		Short: "Builds a dockerfile based on the config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var dckFile []string
			gitignore, err := ignore.CompileIgnoreFile(".gitignore")
			if err != nil {
				return err
			}

			cfg, err := config.Read()
			if err != nil {
				return err
			}

			dckFile = append(dckFile, fmt.Sprintf("FROM %s", cfg.Image.URI))

			dckFile = append(dckFile, getEnvVars(cfg.Image.EnvVars)...)

			// Add precmds before adding any packages
			dckFile = append(dckFile, appendRun(cfg.Image.OnCreation)...)

			// Update base image
			pkgUpdate := pkgm.NewOperation("update", []string{}, []string{}, true)
			update, _ := pkgUpdate.StringCommand(&cfg.Image.Installer)
			dckFile = append(dckFile, fmt.Sprintf("RUN %s", update))

			// Add packages
			pkgInstall := pkgm.NewOperation("add", cfg.Packages, []string{}, true)
			if devBuild {
				pkgInstall.Packages = append(pkgInstall.Packages, cfg.DevPackages...)
			}
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

			// Mounts the current directory to the container's workspace
			if includeFiles {
				dckIgnore := selectDck()
				switch dckIgnore {
				case true:
					dckFile = append(dckFile, fmt.Sprintf("COPY . %s", cfg.Podman.Container.WorkDir))
				case false:
					dckFile = append(dckFile, mountWorkspace(cfg.Podman.Container.WorkDir, gitignore)...)
				}
			}

			if command != "" {
				if _, ok := cfg.Commands[command]; !ok {
					return glg.Errorf("command %s not found in config file", command)
				}

				dckFile = append(dckFile, fmt.Sprintf("ENTRYPOINT [\"%s\", \"-c\", \"%s\"]", cfg.Podman.Container.Shell, cfg.Commands[command]))
			} else {
				dckFile = append(dckFile, fmt.Sprintf("ENTRYPOINT [\"%s\"]", cfg.Podman.Container.Shell))
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
		_, err := file.WriteString(line + "\n\n")

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

// Mounts the current directory to the container's workspace
// and copies any file that doesn't match the .gitignore
func mountWorkspace(workspace string, gitignore *ignore.GitIgnore) []string {
	lines := []string{
		fmt.Sprintf("WORKDIR %s", workspace),
		fmt.Sprintf("VOLUME [\"%s\"]", workspace),
	}
	files, err := container.GetFolderFiles(".")
	if err != nil {
		glg.Fatal(err)
	}

	for _, file := range files {
		if !gitignore.MatchesPath(file) {
			lines = append(lines, fmt.Sprintf("COPY %s %s/%s", file, workspace, file))
		}
	}

	return lines
}

// Returns a list of string that sets the environment variables
// in the Dockerfile
func getEnvVars(vars map[string]string) []string {
	var lines []string
	for key, value := range vars {
		lines = append(lines, fmt.Sprintf("ENV %s=%s", key, value))
	}
	return lines
}

func init() {
	Build.Flags().StringVarP(&command, "command", "c", "", "Command from config file to run on container start")
	Build.Flags().BoolVarP(&includeFiles, "include-files", "i", false, "Include files in current directory in container")
	Build.Flags().BoolVarP(&devBuild, "dev", "d", false, "Include dev packages")
}
