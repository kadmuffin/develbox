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

package config

import (
	"encoding/json"
	"os"
)

func Read() (Struct, error) {
	data, err := os.ReadFile(".develbox/config.json")
	if err != nil {
		return Struct{}, err
	}

	var configs Struct
	err = json.Unmarshal(data, &configs)

	if err != nil {
		return configs, err
	}

	SetDefaults(&configs)
	return configs, nil
}

func WriteConfig(configs *Struct) error {
	os.Mkdir(".develbox", 0755)
	data, _ := json.MarshalIndent(configs, "", "	")

	err := os.WriteFile(".develbox/config.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
