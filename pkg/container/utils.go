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

package container

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

// GetFolderFiles reads all the names in the a path, and returns them
func GetFolderFiles(path string) ([]string, error) {
	// If an item is a folder it will return it like this:
	// "foldername/" instead of "foldername"
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			names = append(names, file.Name()+"/")
		} else {
			names = append(names, file.Name())
		}
	}
	return names, nil
}

// GetFolderFilesMtch returns a list of all the files inside a path that match a string
func GetFolderFilesMtch(path string, match string) ([]string, error) {
	files, err := GetFolderFiles(path)

	if err != nil {
		return []string{}, err
	}

	matches := []string{}
	for _, v := range files {
		if strings.Contains(v, match) {
			matches = append(matches, v)
		}
	}
	return matches, nil
}

// FileExists checks if a file/path exists using os.Stat()
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// processPorts returns a string with all the ports to publish
func processPorts(cfg config.Structure) string {
	return "-p=" + strings.Join(cfg.Container.Ports, "-p=")
}

// processMounts returns a string with the extra volumes to mount
func processMounts(cfg config.Structure) (result []string) {
	for _, v := range cfg.Container.Mounts {
		result = MountArg(result, v, false, "")
	}
	return result
}

// RunCommandList loops through the commands list and runs each one separately
func RunCommandList(name string, commands []string, pman *podman.Podman, root bool, attach podman.Attach) error {
	for _, command := range commands {
		if err := pman.Exec([]string{name, podman.ReplaceEnvVars(command)}, map[string]string{}, true, root, attach).Run(); err != nil {
			return err
		}
	}
	return nil
}

// contains loops through a list to check if a string is inside it
func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// ReplaceEnvVars matches any text that starts with "$" and replaces it with the value of the environment variable. If not set, it will return an empty string.
func ReplaceEnvVars(str string) string {
	re := regexp.MustCompile(`\$(\w+)`)
	// We'll replace any "~" with the home directory and join the path
	if strings.Contains(str, "~/") {
		str = filepath.Join(os.Getenv("HOME"), strings.Replace(str, "~/", "", 1))
		glg.Debugf("Replaced ~ with $HOME path: %s", str)
	}

	return re.ReplaceAllStringFunc(str, func(s string) string {
		envvar := os.Getenv(s[1:])
		if envvar == "" {
			glg.Fatalf("[cfg->SharedFolders] Environment variable '%s' is not set!", s[1:])
		}
		return envvar
	})

}

var mountRegex = regexp.MustCompile(`(?P<host>.*):(?P<container>.*)`)

// MountArg appends a string to a list with the mount argument. It also checks path existance and replaces any environment variables
func MountArg(list []string, mount string, readOnly bool, bindPropagation string) []string {
	// parse %s:%s
	match := mountRegex.FindStringSubmatch(mount)
	if len(match) != 3 {
		glg.Fatalf("Invalid mount format: %s", mount)
	}

	// Save modified paths
	hostPath := ReplaceEnvVars(match[1])
	containerPath := ReplaceEnvVars(match[2])

	// Check if the host path exists
	if !FileExists(hostPath) {
		glg.Warnf("Host path '%s' does not exist! Skipping mount", hostPath)
		return list
	}

	// Set extra options
	extraOpts := ""
	if readOnly {
		extraOpts += ",readonly"
	}
	if bindPropagation != "" {
		extraOpts += fmt.Sprintf(",bind-propagation=%s", bindPropagation)
	}

	return append(list, fmt.Sprintf("--mount=type=bind,src=%s,dst=%s%s", hostPath, containerPath, extraOpts))
}
