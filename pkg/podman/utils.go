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

package podman

import (
	"os"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
)

// Checks if we are inside a container
//
// To do this we check the /run directory for .containerenv (podman)
// or .dockerenv (docker)
func InsideContainer() bool {
	inCtnr := config.FileExists("/run/.containerenv") || config.FileExists("/.dockerenv")

	if os.Getenv("CODESPACES") == "true" {
		inCtnr = false
	}

	return inCtnr
}

// Replaces the following instances:
//
// $USER -> os.Getenv("USER")
//
// $HOME -> /home/$USER
//
// $PWD -> os.Getwd()
func ReplaceEnvVars(s string) string {
	s = strings.ReplaceAll(s, "$$USER", os.Getenv("USER"))
	s = strings.ReplaceAll(s, "$$HOME", "/home/"+os.Getenv("USER"))
	s = strings.ReplaceAll(s, "$$PWD", os.Getenv("PWD"))
	return s
}
