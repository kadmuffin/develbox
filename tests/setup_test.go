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

package main_test

import (
	"os"
	"os/exec"
	"strings"

	"github.com/kadmuffin/develbox/cmd"
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/container"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

var (
	// registryURL is the registry URL to use for the tests
	registryURL string

	// The name of the test container
	testContainerName = "develbox-test"

	// The name of the test image
	testImageName = "alpine:latest"

	podmanPath = config.GetContainerTool()
	pman       = podman.New(podmanPath)

	setupAlreadyRun = false

	// keepContainer is a flag to keep the container after the tests are done.
	// Only works outside of GitHub Actions (for some reason)
	keepContainer = false
)

// Setup sets up the test environment
func Setup(copyCfg bool, createContainer bool) {
	// Set up the logger
	glg.Get().SetMode(glg.STD).SetLevel(glg.DEBG).SetWriter(os.Stdout)

	// Make sure the test environment is clean
	// This is to make sure that the tests don't interfere with each other
	// and that the tests don't leave any files behind
	err := cleanTestEnv()
	if err != nil {
		glg.Fatalf("Failed to clean test environment: %s", err)
	}

	PullImage()

	// Copy the config file
	if copyCfg {
		CopyConfig()
	}

	if keepContainer {
		// Check if the container already exists
		exists := ContainerExists(testContainerName)
		switch exists {
		case false:
			createContainer = true
		case true:
			createContainer = false
			if os.Getenv("GITHUB_ACTIONS") != "true" {
				_, err := exec.Command(podmanPath, "start", testContainerName).CombinedOutput()
				if err != nil {
					glg.Fatalf("Failed to start container: %s", err)
				}
			}
		}
	}

	// Create the test container
	if createContainer {

		// Create a container
		container.PkgVersion = cmd.GetRootCLI().Version
		container.DontStopOnFinish = true
		err := container.Create(SampleConfig, true)
		if err != nil {
			glg.Fatalf("Failed to create container: %s", err)
		}

		// Check if the container exists
		exists := ContainerExists(testContainerName)

		if !exists {
			glg.Fatalf("Container %s does not exist", testContainerName)
		}
	}

	keepContainer = false
	setupAlreadyRun = true
}

// Clean cleans the test environment
func cleanTestEnv() error {
	// Delete and create a new .develbox folder
	if !keepContainer {
		os.RemoveAll(".develbox")
		os.MkdirAll(".develbox/home", 0755)
	}

	// Get the list of containers
	out, err := exec.Command(podmanPath, "ps", "-a", "-q").Output()
	if err != nil {
		return glg.Errorf("Failed to get list of containers: %s", err)
	}

	// Remove all containers with the test name
	if len(out) > 0 && !keepContainer {
		// Remove the trailing newline
		out = out[:len(out)-1]

		// Split the output into a list of containers
		containers := strings.Split(string(out), "\n")

		// Remove each container if it has the name ${testContainerName}
		for _, container := range containers {
			if strings.HasPrefix(container, testContainerName) {
				_, err := exec.Command(podmanPath, "rm", "-f", container).Output()
				if err != nil {
					return glg.Errorf("Failed to remove container %s: %s", container, err)
				}

				glg.Debugf("Removed container %s", container)
			}
		}

	}

	return nil
}

// CopyConfig copies the config file to .develbox/config.json
func CopyConfig() {
	os.MkdirAll(".develbox", 0755)
	// Copy the config file
	out, err := exec.Command("cp", "config/alpine.json", ".develbox/config.json").CombinedOutput()
	if err != nil {
		glg.Fatalf("Failed to copy config file: %s", string(out))
	}
}
