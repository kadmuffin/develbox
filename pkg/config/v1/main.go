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

// Package v1_config defines the configuration file and provides functions to read and write it
//
// This is the old version of the config file, which is now deprecated in favor of the new version (v2 at time of writing)
package v1_config

import (
	"encoding/json"
	"os"

	"github.com/kpango/glg"
)

// Read reads the config file and returns the Struct
func Read() (Struct, error) {
	return ReadFile(".develbox/config.json")
}

// ReadFile reads the config file  from a path and returns the Struct
func ReadFile(path string) (Struct, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Struct{}, err
	}

	glg.Infof("Reading config file at %s", path)
	return ReadBytes(data)
}

// Write writes the config file
func Write(configs *Struct) error {
	os.Mkdir(".develbox", 0755)
	data, _ := json.MarshalIndent(configs, "", "	")

	err := os.WriteFile(".develbox/config.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads the config file from bytes and returns the Struct
func ReadBytes(data []byte) (Struct, error) {
	var configs Struct
	err := json.Unmarshal(data, &configs)

	if err != nil {
		return configs, err
	}

	SetName(&configs)
	return configs, nil
}
