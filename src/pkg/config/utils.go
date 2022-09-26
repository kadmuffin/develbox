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

package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/kpango/glg"
)

// Checks whether the config folder exists.
//
// Wrapper around config.FileExists()
func ConfigFolderExists() bool {
	return FileExists(".develbox")
}

// Checks whether the config file exists.
//
// Wrapper around config.FileExists()
func ConfigExists() bool {
	return FileExists(".develbox/config.json")
}

// Gets the current folder's full path.
//
// Throws a fatal error in case of failure
// and exists the program.
func GetCurrentDirectory() string {
	currentDir, err := os.Getwd()

	if err != nil {
		glg.Fatalf("failed to get current directory:\n	%s", err)
	}

	return currentDir
}

// Checks if a file/path exists.
// Returns true if it exists
//
// Wrapper around os.Stat() & os.IsNotExists()
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns a hash string made using the current directory's name.
func GetDirNmHash() string {
	currentDirName := filepath.Base(GetCurrentDirectory())
	hasher := sha256.New()
	hasher.Write([]byte(currentDirName))
	dir := hasher.Sum(nil)
	return hex.EncodeToString(dir)

}
