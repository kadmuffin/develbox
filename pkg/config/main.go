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

// Package config auto detects the version of the config file and returns the Struct in the latest version
//
// It converts the v1 config file to the v2 config file if it detects a v1 config file
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	v1_config "github.com/kadmuffin/develbox/pkg/config/v1"
	"github.com/kpango/glg"
)

// Read reads the config file and returns the Struct
func Read() (cfg Struct, err error) {
	var v1Cfg bool
	cfg, err, v1Cfg = ReadFile(".develbox/config.json")
	if err == nil && v1Cfg {
		err = WriteNewVersion(&cfg)
	}

	return cfg, err
}

// ReadFile reads the config file  from a path and returns the Struct
//
// It converts the v1 config file to the v2 config file if it detects a v1 config file
func ReadFile(path string) (Struct, error, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Struct{}, err, false
	}

	glg.Infof("Reading config file at %s", path)

	configs, err, v1Cfg := ReadBytes(data)
	if err != nil {
		return Struct{}, err, v1Cfg
	}

	return configs, nil, v1Cfg
}

// Write writes the config file
func Write(configs *Struct) error {
	os.Mkdir(".develbox", 0755)
	data, _ := json.MarshalIndent(configs, "", "  ")

	err := os.WriteFile(".develbox/config.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads the config file from bytes and returns the Struct
//
// It converts the v1 config file to the v2 config file if it detects a v1 config file
func ReadBytes(data []byte) (Struct, error, bool) {
	var configs Struct
	err := json.Unmarshal(data, &configs)
	if err != nil {
		// If the config file is in the old format, convert it to the new format
		var oldConfigs v1_config.Struct
		err = json.Unmarshal(data, &oldConfigs)
		if err != nil {
			return Struct{}, err, false
		}

		glg.Info("Valid v1 config! Converting v1 config file to new format (v2)...")
		configs = ConvertFromV1(oldConfigs)
		return configs, nil, true
	}

	return configs, nil, false
}

// Exists checks if the config file exists
func Exists() bool {
	_, err := os.Stat(".develbox/config.json")
	return err != nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

// GetCurrentDirectory returns the current directory
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
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

// WriteNewVersion writes the new version of the config file
func WriteNewVersion(configs *Struct) error {
	glg.Warn("Updating v1 config file to new format (v2)... (a backup of the old config file will be saved as config.json.bak)")

	// Backup the old config file
	err := os.Rename(".develbox/config.json", ".develbox/config.json.bak")
	if err != nil {
		return err
	}

	// Write the new config file
	return Write(configs)
}
