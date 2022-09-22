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

package develbox

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func parseJson(bytes []byte) DevSetings {
	var configs DevSetings
	err := json.Unmarshal(bytes, &configs)

	if err != nil {
		log.Fatalf("Couldn't parse the config file, exited with: %s", err)
	}

	SetContainerName(&configs)

	return configs
}

func SetContainerName(config *DevSetings) {
	if config.Podman.Container.Name == "" {
		hasher := sha256.New()
		hasher.Write([]byte(filepath.Base(getCurrentDirectory())))
		dir := hasher.Sum(nil)
		config.Podman.Container.Name = hex.EncodeToString(dir)
	}
}

func ReadConfig() DevSetings {
	data, err := os.ReadFile(".develbox/config.json")
	if err != nil {
		log.Fatalf("Couldn't read the file .develbox/config.json, exited with: %s", err)
	}

	configs := parseJson(data)
	SetContainerName(&configs)

	return configs
}

func WriteConfig(configs *DevSetings) {
	data, _ := json.MarshalIndent(configs, "", "	")

	err := os.WriteFile(".develbox/config.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func ConfigFolderExists() bool {
	_, err := os.Stat(".develbox")

	return !os.IsNotExist(err)
}

func ConfigExists() bool {
	_, err := os.Stat(".develbox/config.json")

	return !os.IsNotExist(err)
}
