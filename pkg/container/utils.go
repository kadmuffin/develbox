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
	"os"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
)

// Reads all the names in the a path, and returns them
// If an item is a folder it will return it like this:
// "foldername/" instead of "foldername"
func GetFolderFiles(path string) ([]string, error) {
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

// Returns a list of all the files inside a path that match a string
//
// Wrapper around GetFolderFiles() that filters the results
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

// Checks if a file/path exists using os.Stat()
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns a string with all the ports to publish
func processPorts(cfg config.Struct) string {
	return "-p=" + strings.Join(cfg.Podman.Container.Ports, "-p=")
}

// Returns a string with the extra volumes to mount
func processVolumes(cfg config.Struct) string {
	return "-v=" + strings.Join(cfg.Podman.Container.Mounts, "-v=")
}

// Loops through the commands list and runs each one separately
func RunCommandList(name string, commands []string, pman *podman.Podman, root bool, attach podman.Attach) error {
	for _, command := range commands {
		err := pman.Exec([]string{name, podman.ReplaceEnvVars(command)}, map[string]string{}, true, root, attach).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// Loops through a list to check if a string is inside it
func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
