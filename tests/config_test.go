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

package main_test

import (
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/config/v1config"
)

var (
	SampleConfig = config.Structure{
		Image: config.Image{
			URI:        "alpine:edge",
			OnCreation: []string{},
			OnFinish:   []string{},
			Variables:  map[string]string{},
			PkgManager: config.PackageManager{
				Operations: config.Operations{
					Add:   "apk add {args}",
					Del:   "apk del {args}",
					Upd:   "apk update {args}",
					Upg:   "apk upgrade {args}",
					Srch:  "apk search {args}",
					Clean: "rm -rf /var/cache/apk",
				},
				Modifiers: map[string]string{},
			},
		},
		Podman: config.Podman{
			Args:       []string{"--net=host"},
			Privileged: true,
			Path:       podmanPath,
			Rootless:   true,
			AutoDelete: false,
			AutoCommit: false,
		},
		Container: config.Container{
			Name:     "develbox-test",
			WorkDir:  "/code",
			Shell:    "/bin/sh",
			RootUser: false,
			Binds: config.Binds{
				XOrg:      false,
				Dev:       false,
				Variables: []string{},
			},
			Ports:  []string{},
			Mounts: []string{},
			SharedFolders: map[string]interface{}{
				"alpine": "/var/cache/apk/",
			},
		},
		Commands: map[string]interface{}{
			"test":  "echo test",
			"test2": "!test",
			"test3": "#echo test",
		},
		Packages: []string{
			"nodejs",
			"npm",
		},
		DevPackages: []string{
			"git",
			"make",
		},
		UserPkgs: v1config.UserPkgs{
			Packages:    []string{},
			DevPackages: []string{},
		},
		Experiments: v1config.Experiments{
			Socket: false,
		},
	}
)

func init() {
	// Set the podman path
	config.CheckDocker(&SampleConfig)
}

// TestConfig tests the read related functions of the config package
func TestConfig(t *testing.T) {
	t.Logf("Current location is %s", os.Getenv("PWD"))
	// Run setup steps
	Setup(true, false)

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// Compare the two configs
	if !CompareConfigs(cfg, SampleConfig) {
		t.Fatalf("[config.Read()] Config file does not match template")
	}

	// Now we will test config.ReadBytes()
	// Lets the config using os.ReadFile so we can get
	// the bytes
	bytes, err := os.ReadFile(".develbox/config.json")

	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// Read the config file
	cfg, wasV1Conf, err := config.ReadBytes(bytes)

	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// We expected a v2 config, so we should get false
	if wasV1Conf {
		t.Fatalf("[config.ReadBytes()] Detected v1 config instead of v2")
	}

	// It should be the same as the original config
	if !CompareConfigs(cfg, SampleConfig) {
		t.Fatalf("[config.ReadBytes()] Config file does not match template")
	}
}

// TestConversion tests the conversion of a v1 config to a v2 config
func TestConversion(t *testing.T) {
	// Run setup steps
	Setup(false, false)

	// Copy the v1 config to the .develbox/config.json file
	exec.Command("cp", "config/alpine.v1.json", ".develbox/config.json").Run()

	bytes, err := os.ReadFile(".develbox/config.json")

	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// Read the config file
	v1cfg, wasV1Conf, err := config.ReadBytes(bytes)

	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// We expected a v1 config, so we should get true
	if !wasV1Conf {
		t.Fatalf("[config.Read() V1Cfg] Detected v2 config instead of v1")
	}

	// It should be aproximately the same as the original config
	if !CompareConfigs(v1cfg, SampleConfig) {
		t.Fatalf("[config.Read() V1Cfg] Config file does not match template")
	}
}

// TestWriteConfig tests writing a config template
func TestWriteConfig(t *testing.T) {
	// Run setup steps
	Setup(false, false)

	err := config.Write(&SampleConfig)
	if err != nil {
		t.Fatalf("Failed to write config file: %s", err)
	}

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		t.Fatalf("Failed to read config file: %s", err)
	}

	// Compare the two configs
	if !CompareConfigs(cfg, SampleConfig) {
		t.Fatalf("Config file does not match template")
	}
}

// CompareConfigs compares two config.Structure structs
func CompareConfigs(cfg1, cfg2 config.Structure) bool {
	// Compare the two configs
	return reflect.DeepEqual(cfg1, cfg2)
}
