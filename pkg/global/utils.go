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

// Package global contains global variables and functions
package global

import (
	"os"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
)

// GetDataHome gets XDG_DATA_HOME, if not set, use ~/.local/share and set it
func GetDataHome() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = os.Getenv("HOME") + "/.local/share"
		os.Setenv("XDG_DATA_HOME", dataHome)
	}
	return dataHome
}

// GetConfigHome gets XDG_CONFIG_HOME, if not set, use ~/.config and set it
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

// GetLastPathPart get last part of a path. Example: /home/user/Downloads -> Downloads
func GetLastPathPart(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

// GetPathBeforeLastPart gets everything before the last part of a path. Example: /home/user/Downloads -> /home/user
func GetPathBeforeLastPart(path string) string {
	return path[:strings.LastIndex(path, "/")]
}

// CreateFolder creates a new folder at the given path
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsFolder checks if a path is a folder or a file.
func IsFolder(path string) bool {
	// For doing this, it checks if the last part of the path
	// if it has a "/", it's a folder, otherwise it's a file
	return strings.HasSuffix(path, "/")
}

// HashPath returns a hash of a path and mantains the "/" at the end of the path if it exists.
func HashPath(path string) string {
	hashPath := config.GetPathHash(path)
	if IsFolder(path) {
		hashPath += "/"
	}
	return hashPath
}
