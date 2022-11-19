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
	"testing"

	"github.com/kadmuffin/develbox/cmd"
	"github.com/kadmuffin/develbox/pkg/container"
)

// TestCreate tests the create function
func TestCreate(t *testing.T) {
	Setup(false, false)

	// Create a container
	container.PkgVersion = cmd.GetRootCLI().Version
	err := container.Create(SampleConfig, true)
	if err != nil {
		t.Errorf("Failed to create container: %s", err)
	}

	// Check if the container exists
	exists := pman.Exists(testContainerName)

	// Log system info (if running in github actions)
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		out, _ := exec.Command("uname", "-a").CombinedOutput()
		// Debug stuff
		t.Log(string(out))

		out, _ = exec.Command("docker", "version").CombinedOutput()
		// Debug stuff
		t.Log(string(out))

		out, _ = exec.Command("docker", "info").CombinedOutput()
		// Debug stuff
		t.Log(string(out))
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		out, _ := exec.Command("docker", "inspect", testContainerName).CombinedOutput()
		// Debug stuff
		t.Log(string(out))

		out, _ = exec.Command("docker", "logs", testContainerName).CombinedOutput()
		// Debug stuff
		t.Log(string(out))
	}

	if !exists {
		t.Fatalf("Container %s does not exist", testContainerName)
	}
}

// TestCreateCmd tests the create command (cli)
func TestCreateCmd(t *testing.T) {
	Setup(true, false)

	// Create a container
	container.PkgVersion = cmd.GetRootCLI().Version
	cmd.GetRootCLI().SetArgs([]string{"create", "-f"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Failed to create container: %s", err)
	}

	// Check if the container exists
	exists := pman.Exists(testContainerName)

	if !exists {
		t.Fatalf("Container %s does not exist", testContainerName)
	}
}
