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

package create

import (
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
)

// downloadConfig downloads the config file from the given URL
func downloadConfig(argum string, url string) (config.Structure, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s.json", url, argum))

	if err != nil {
		return config.Structure{}, glg.Fail(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return config.Structure{}, glg.Errorf("Response from source returned bad status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return config.Structure{}, glg.Errorf("Something went wrong while downloading the config file. %s", err)
	}

	cfg, v1Cfg, err := config.ReadBytes(data)
	if err != nil {
		return config.Structure{}, glg.Errorf("Failed to parse the JSON data. %s", err)
	}

	if v1Cfg {
		glg.Warn("The config file is in the old format. Develbox will update it to the new format.")
	}
	return cfg, nil
}

// getKeys returns the keys of a map as a slice
func getKeys(data map[string][]string) []string {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
