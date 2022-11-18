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
	"testing"

	"github.com/kadmuffin/develbox/cmd"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

// TestCreate tests the create function
func TestCreate(t *testing.T) {
	Setup(false)

	// Create a container
	container.PkgVersion = cmd.GetRootCLI().Version
	err := container.Create(SampleConfig, true)
	if err != nil {
		t.Errorf("Failed to create container: %s", err)
	}

	// Check if the container exists
	exists := pman.Exists(testContainerName)

	if !exists {
		t.Fatalf("Container %s does not exist", testContainerName)
	}

	// Remove the container
	err = pman.Remove([]string{testContainerName}, podman.Attach{})
	if err != nil {
		glg.Fatalf("Failed to remove container: %s", err)
	}

	// Check if the container exists
	exists = pman.Exists(testContainerName)

	if exists {
		t.Fatalf("Container %s still exists", testContainerName)
	}
}

// TestCreateCmd tests the create command (cli)
func TestCreateCmd(t *testing.T) {
	Setup(true)

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

	// Remove the container
	err = pman.Remove([]string{testContainerName}, podman.Attach{})

	if err != nil {
		glg.Fatalf("Failed to remove container: %s", err)
	}

	// Check if the container exists
	exists = pman.Exists(testContainerName)

	if exists {
		t.Fatalf("Container %s still exists", testContainerName)
	}
}
