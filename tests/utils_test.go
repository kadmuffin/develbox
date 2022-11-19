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

// This file contains utility functions for tests
package main_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

// Detect if we are running under GitHub Actions
func IsGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// GetRegistryURL gets the registry url
// By default, it will use docker.io for the registry
// If running under GitHub Actions, it will use the GitHub Container Registry
func GetRegistryURL() string {
	// Disabled until I figure this out
	//if IsGitHubActions() {
	//	return "ghcr.io"
	//}

	return "docker.io"
}

// PullImage pulls the latest image
func PullImage() {
	if !setupAlreadyRun {
		registryURL = GetRegistryURL()
		image := fmt.Sprintf("%s/%s", registryURL, testImageName)

		// Pull the latest image
		glg.Infof("Pulling image %s", image)
		out, err := exec.Command(podmanPath, "pull", image).CombinedOutput()
		if err != nil {
			glg.Fatalf("Failed to pull image: %s", string(out))
		}
	}
}

// ContainerExists checks if a container exists
func ContainerExists(name string) bool {
	// Check if the container exists

	switch pman.IsDocker() {
	case true:
		cmd := pman.RawCommand([]string{"inspect", name}, podman.Attach{})

		_, err := cmd.Output()
		if err != nil {
			glg.Errorf("Container %s does not exist", name)
			glg.Debug(err)
			return false
		}

	case false:
		cmd := pman.RawCommand([]string{"container", "exists", name}, podman.Attach{})
		_, err := cmd.Output()
		if err != nil {
			glg.Errorf("Container %s does not exist", name)
			glg.Debug(err)
			return false
		}

	}
	return true
}
