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

package global

import (
	"os"
	"strings"
)

// Get XDG_DATA_HOME, if not set, use ~/.local/share and set it
// In the future, this will be used to store the global folders
// for the containers
func GetDataHome() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = os.Getenv("HOME") + "/.local/share"
		os.Setenv("XDG_DATA_HOME", dataHome)
	}
	return dataHome
}

// Get XDG_CONFIG_HOME, if not set, use ~/.config and set it
func GetConfigHome() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = os.Getenv("HOME") + "/.config"
		os.Setenv("XDG_CONFIG_HOME", configHome)
	}
	return configHome
}

// Get XDG_CACHE_HOME, if not set, use ~/.cache and set it
func GetCacheHome() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = os.Getenv("HOME") + "/.cache"
		os.Setenv("XDG_CACHE_HOME", cacheHome)
	}
	return cacheHome
}

// Get last part of a path
// Example: /home/user/Downloads -> Downloads
func GetLastPathPart(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

// Create a new folder at the given path
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}
