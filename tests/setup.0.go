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

// This files sets up the test environment for the tests in this directory.
// And provides some helper functions for the tests.

package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/kpango/glg"
)

var (
	// The registry URL to use for the tests
	registryURL string

	// The name of the test container
	testContainerName = "develbox-test"

	// The name of the test image
	testImageName = "alpine:latest"
)

// Setup the test environment
func Setup(t *testing.T) {
	// Set up the logger
	glg.Get().SetMode(glg.STD).SetLevel(glg.DEBG).SetWriter(os.Stdout)

	// Make sure the test environment is clean
	// This is to make sure that the tests don't interfere with each other
	// and that the tests don't leave any files behind
	err := cleanTestEnv()
	if err != nil {
		t.Fatalf("Failed to clean test environment: %s", err)
	}

	PullImage(t)
}

// Clean the test environment
func cleanTestEnv() error {
	// Get the list of containers
	out, err := exec.Command("podman", "ps", "-a", "-q").Output()
	if err != nil {
		return glg.Errorf("Failed to get list of containers: %s", err)
	}

	// Remove all containers
	if len(out) > 0 {
		// Remove the trailing newline
		out = out[:len(out)-1]

		// Split the output into a list of containers
		containers := strings.Split(string(out), "\n")

		// Remove each container if it has the name ${testContainerName}
		for _, container := range containers {
			if strings.HasPrefix(container, testContainerName) {
				_, err := exec.Command("podman", "rm", "-f", container).Output()
				if err != nil {
					return glg.Errorf("Failed to remove container %s: %s", container, err)
				}

				glg.Debugf("Removed container %s", container)
			}
		}

	}

	return nil
}
