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

	"github.com/kadmuffin/develbox/pkg/container"
)

// TestEnter tests the enter function
func TestEnter(t *testing.T) {
	keepContainer = true
	Setup(false, true)
	// Enter the container
	container.DontAttachEnter = true
	err := container.Enter(SampleConfig, false)
	if err != nil {
		t.Errorf("Failed to enter container: %s", err)
	}
}

// TestEnterInstall tests the InstallAndEnter function
func TestEnterInstall(t *testing.T) {
	keepContainer = true
	Setup(false, true)
	// Enter the container
	container.DontAttachEnter = true
	err := container.InstallAndEnter(SampleConfig, true)
	if err != nil {
		t.Errorf("Failed to enter container: %s", err)
	}
}
