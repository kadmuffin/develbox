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
	"strings"

	"github.com/kpango/glg"
)

type Podman struct {
	path string
}

type Attach struct {
	Stdin     bool
	Stdout    bool
	Stderr    bool
	PseudoTTY bool
}

func New(path string) Podman {
	glg.Debugf("Creating podman instance using '%s'.", path)
	if err := exec.Command(path, "version").Run(); err != nil {
		glg.Fatalf("Can't access the podman executable: %s", err)
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
func (e *Podman) Create(args []string, attach Attach) *exec.Cmd {
	params := []string{"run", "-d", "-t", "--init"}
	params = append(params, args...)

	cmd := e.cmd(params, attach)
	glg.Debugf("Creating container using the following command: %s", cmd.String())

	return cmd
}

// Executes a command inside a running container and
// attaches (Stdin, Stdout) if "attach" is true.
func (e *Podman) Exec(args []string, sh bool, root bool, attach Attach) *exec.Cmd {
	params := []string{"exec", "-i"}

	if attach.PseudoTTY {
		params = append(params, "-t")
	}

	if attach == *new(Attach) {
		params = append(params, "-d")
	}

	if root {
		params = append(params, "--user", "0:0")
	}

	params = append(params, args[0])

	if sh {
		params = append(params, "sh", "-c")
	}

	params = append(params, args[1:]...)

	cmd := e.cmd(params, attach)
	glg.Debugf("Executing command: %s", cmd.String())

	return cmd
}

// Returns a boolean that indicates if the container was found.
func (e *Podman) Exists(name string) bool {
	params := []string{"container", "exists", name}

	// Docker doesn't have an exists function, so this is
	// the closest thing I could find.
	if strings.Contains(e.path, "docker") {
		params = []string{"ps", "-a", "|", "grep", name}
	}

	return e.cmd(params, Attach{}).Run() == nil
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
	e.Stop(args, attach)

	params := []string{"rm"}
	params = append(params, args...)

	glg.Debugf("Removing container using the following arguments:\n  - %s", params)

	return e.cmd(params, attach).Run()
}

// Copies files into the container
func (e *Podman) Copy(args []string, attach Attach) *exec.Cmd {
	e.Stop(args, attach)

	params := []string{"cp"}
	params = append(params, args...)

	glg.Debugf("Running copy using the following arguments:\n  - %s", params)

	return e.cmd(params, attach)
}

// Checks if the path contains the word "docker"
//
// Mainly here so we can make this work on docker too.
func (e *Podman) IsDocker() bool {
	return strings.Contains(e.path, "docker")
}
