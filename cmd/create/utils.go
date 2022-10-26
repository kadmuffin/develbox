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
	"os"
	"sort"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
)

func downloadConfig(argum string, url string) config.Struct {
	resp, err := http.Get(fmt.Sprintf("%s/%s.json", url, argum))

	if err != nil {
		glg.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		glg.Errorf("Response from source returned bad status: %s", resp.StatusCode)
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Fatalf("Something went wrong while downloading the config file. %s", err)
	}

	cfg, err := config.ReadBytes(data)
	if err != nil {
		glg.Fatalf("Failed to parse the JSON data. %s", err)
	}
	return cfg
}

// Returns the keys of a map as a slice
func getKeys(data map[string][]string) []string {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
