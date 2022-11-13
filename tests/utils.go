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
package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// Detect if we are running under GitHub Actions
func IsGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// Gets the registry url
// By default, it will use the local registry
// If running under GitHub Actions, it will use the GitHub Container Registry
func GetRegistryURL() string {
	if IsGitHubActions() {
		return "ghcr.io"
	}

	return "localhost:5000"
}

// Pull the latest image
func PullImage(t *testing.T) {
	registryURL = GetRegistryURL()
	image := fmt.Sprintf("%s/%s", registryURL, testImageName)

	// Pull the latest image
	out, err := exec.Command("podman", "pull", image).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to pull image: %s", string(out))
	}
}
