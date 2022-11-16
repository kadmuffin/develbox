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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	v1config "github.com/kadmuffin/develbox/pkg/config/v1config"
	"github.com/kpango/glg"
	"github.com/spf13/viper"
)

// Read reads the config file and returns the Struct
func Read() (cfg Structure, err error) {
	var v1Cfg bool
	cfg, v1Cfg, err = ReadFile(".develbox/config.json")
	if err == nil && v1Cfg {
		err = WriteNewVersion(&cfg)
	}

	return cfg, err
}

// ReadFile reads the config file  from a path and returns the Struct
//
// It converts the v1 config file to the v2 config file if it detects a v1 config file
func ReadFile(path string) (Structure, bool, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return Structure{}, false, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Structure{}, false, err
	}
	defer file.Close()

	return parseWithViper(file)
}

// ReadBytes parses bytes and returns the Struct
//
// It converts the v1 config file to the v2 config file if it detects a v1 config file
func ReadBytes(data []byte) (parsed Structure, wasV1Conf bool, err error) {
	buffer := bytes.NewBuffer(data)
	viper.SetConfigType("json")

	err = viper.ReadConfig(buffer)
	if err != nil {
		return Structure{}, false, err
	}

	return parseWithViper(buffer)
}

// Write writes the config file
func Write(configs *Structure) error {
	glg.Infof("Writing config file to %s", ".develbox/config.json")

	err := os.MkdirAll(".develbox", 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(".develbox/config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(configs)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if the config file exists
func Exists() bool {
	_, err := os.Stat(".develbox/config.json")
	exists := err == nil
	if exists {
		glg.Info("Config file exists!")
	} else {
		glg.Info("Config file does not exist!")
	}
	return exists
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
func WriteNewVersion(configs *Structure) error {
	glg.Warn("Updating v1 config file to new format (v2)... (a backup of the old config file will be saved as config.json.bak)")

	// Backup the old config file
	err := os.Rename(".develbox/config.json", ".develbox/config.json.bak")
	if err != nil {
		return err
	}

	// Write the new config file
	return Write(configs)
}

// parseWithViper assumes viper is already configured and returns the parsed config
func parseWithViper(reader io.Reader) (Structure, bool, error) {
	// Use json parser until I can figure out how to use the viper parser
	// properly (the issue arises from parsing interface{} types, specifically, shared-folders, see pkg/config/v1/struct.go:125)
	decoder := json.NewDecoder(reader)

	if !viper.IsSet("container") && viper.IsSet("podman.container") {
		var v1Struct v1config.Struct

		err := decoder.Decode(&v1Struct)
		if err != nil {
			return Structure{}, false, err
		}

		parsed := ConvertFromV1(&v1Struct)
		return parsed, true, nil
	}

	var parsed Structure
	err := decoder.Decode(&parsed)
	if err != nil {
		return Structure{}, false, nil
	}

	SetDefaults(&parsed)

	return parsed, false, nil
}
