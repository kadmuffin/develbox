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
	"os/exec"

	"github.com/kpango/glg"
)

type Podman struct {
	path    string
	lastCmd exec.Cmd
}

type Attach struct {
	Stdin  bool
	Stdout bool
	Stderr bool
}

func New(path string) Podman {
	glg.Debugf("Creating podman instance using '%s'.", path)
	if err := exec.Command("podman").Run(); err != nil {
		glg.Fatalf("Can't access the podman executable: %w", err)
	}
	return Podman{path: path}
}

// Private function that manages the command creation.
// Created as a boilerplate for other public functions,
func (e *Podman) cmd(args []string, attach Attach) *exec.Cmd {
	cmd := exec.Command(e.path, args...)

	if attach.Stdin {
		cmd.Stdin = os.Stdin
	}
	if attach.Stdout {
		cmd.Stdout = os.Stdout
	}
	if attach.Stderr {
		cmd.Stderr = os.Stderr
	}

	return cmd
}

// Creates a container using the arguments provided.
func (e *Podman) Create(args []string, attach Attach) error {
	params := []string{"run", "-d", "-t", "--init"}
	params = append(params, args...)

	glg.Debugf("Creating container using the following arguments:\n  - %s", params)

	return e.cmd(params, attach).Run()
}

// Executes a command inside a running container and
// attaches (Stdin, Stdout) if "attach" is true.
func (e *Podman) Exec(args []string, sh bool, root bool, attach Attach) error {
	params := []string{"exec"}
	if root {
		params = append(params, "--user", "0:0")
	}

	params = append(params, "sh", "-c")
	params = append(params, args...)

	glg.Debugf("Executing a command using %s:\n  - Command: %s", params)

	return e.cmd(params, attach).Run()
}

// Returns a boolean that indicates if the container was found.
func (e *Podman) Exists(name string) bool {
	params := []string{"container", "exists"}
	params = append(params, name)

	return e.cmd(params, Attach{}).Run() != nil
}

// Starts a container and returns an error in case of failure.
//
// In arguments, the first argument should probably be the
// container's name/id.
func (e *Podman) Start(args []string, attach Attach) error {
	params := []string{"start"}
	params = append(params, args...)

	glg.Debugf("Starting container using the following arguments:\n  - %s", params)

	return e.cmd(params, attach).Run()
}

// Stops a container and returns an error in case of failure.
//
// In arguments, the first argument can should the
// container's name/id if no flag are added before of the name.
func (e *Podman) Stop(args []string, attach Attach) error {
	params := []string{"stop"}
	params = append(params, args...)

	glg.Debugf("Stopping container using the following arguments:\n  - %s", params)

	return e.cmd(params, attach).Run()
}

// Removes a container and returns an error in case of failure.
//
// In arguments, the first argument can should the
// container's name/id if no flag are added before of the name.
func (e *Podman) Remove(args []string, attach Attach) error {
	params := []string{"rm"}
	params = append(params, args...)

	glg.Debugf("Removing container using the following arguments:\n  - %s", params)

	return e.cmd(params, attach).Run()
}

// Returns the last ran command
//
// Useful if you need to access it's output or other information.
func (e *Podman) GetLastCommand() exec.Cmd {
	return e.lastCmd
}
