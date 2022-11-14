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

// Package config has the config file struct. This file has utils for the config file
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/kpango/glg"
)

// ConfigFolderExists returns true if the config folder exists.
func ConfigFolderExists() bool {
	return FileExists(".develbox")
}

// ConfigExists returns true if the config file exists.
func ConfigExists() bool {
	return FileExists(".develbox/config.json")
}

// GetCurrentDirectory returns the current folder's full path.
func GetCurrentDirectory() string {
	currentDir, err := os.Getwd()

	if err != nil {
		glg.Fatalf("failed to get current directory:\n	%s", err)
	}

	return currentDir
}

// FileExists returns true if a file/path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetDirNmHash returns a hash made using the current directory's name.
func GetDirNmHash() string {
	currentDirName := filepath.Base(GetCurrentDirectory())
	hasher := sha256.New()
	hasher.Write([]byte(currentDirName))
	dir := hasher.Sum(nil)
	return hex.EncodeToString(dir)
}

// GetPathHash returns a hash made using the provided path.
func GetPathHash(path string) string {
	hasher := sha256.New()
	hasher.Write([]byte(path))
	dir := hasher.Sum(nil)
	return hex.EncodeToString(dir)
}
