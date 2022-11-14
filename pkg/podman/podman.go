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

// Package podman is a wrapper around os/exec to run podman commands.
package podman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/kpango/glg"
)

// Podman is a struct that saves the path to the podman executable.
type Podman struct {
	path string
}

// Attach is config struct that sets the Stdin, Stdout and Stderr
type Attach struct {
	// Stdin sets if stdin should be attached to the current process.
	Stdin bool
	// Stdout sets if stdout will be attached to the current process.
	Stdout bool
	// Stderr sets if stderr will be attached to the current process.
	Stderr bool
	// PseudoTTY sets if podman should allocate a pseudo-TTY for the container.
	PseudoTTY bool
}

// New creates a new Podman struct with the path to the podman executable.
func New(path string) Podman {
	glg.Infof("Podman path set to '%s'.", path)
	cmd := exec.Command(path, "--version")
	glg.Infof("Verifying that podman exists using '%s'.", cmd.String())

	if err := cmd.Run(); err != nil {
		glg.Fatalf("Can't access the podman executable: %s", err)
	}

	return Podman{path: path}
}

// cmd is private function that manages the command creation. Created as a boilerplate for other public functions.
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

// Create creates a container using the arguments provided.
func (e *Podman) Create(args []string, attach Attach) *exec.Cmd {
	params := []string{"run", "-t", "--init"}
	params = append(params, args...)

	return PrintCommandR("Creating container using the following command: %s", e.cmd(params, attach))

}

// Exec executes a command inside a running container and attaches (Stdin, Stdout) if "attach" is true.
func (e *Podman) Exec(args []string, envVars map[string]string, sh bool, root bool, attach Attach) *exec.Cmd {
	uid := os.Getuid()
	params := []string{"exec", "-i"}

	if attach.PseudoTTY {
		params = append(params, "-t")
	}

	if attach == *new(Attach) {
		params = append(params, "-d")
	}

	for k, v := range envVars {
		params = append(params, "-e", fmt.Sprintf("%s=%s", k, ReplaceEnvVars(v)))
	}

	if root {
		params = append(params, "--user", "0:0")
	} else {
		params = append(params, "--user", fmt.Sprintf("%d:%d", uid, uid))
	}

	params = append(params, args[0])

	if sh {
		params = append(params, "sh", "-c")
	}

	params = append(params, args[1:]...)

	return PrintCommandR("Executing command: %s", e.cmd(params, attach))
}

// Exists returns a boolean that indicates if the container was found.
func (e *Podman) Exists(name string) bool {
	params := []string{"container", "exists", name}

	// Docker doesn't have an exists function, so this is
	// the closest thing I could find.
	if e.IsDocker() {
		params = []string{"inspect", name}
	}

	return e.cmd(params, Attach{}).Run() == nil
}

// Start starts a container and returns an error in case of failure. The first argument has to be the container's name/id.
func (e *Podman) Start(args []string, attach Attach) error {
	params := []string{"start"}
	params = append(params, args...)

	return PrintCommandR("Starting container using the following arguments:\n  - %s", e.cmd(params, attach)).Run()
}

// Stop stops a container and returns an error in case of failure. In arguments, the first argument has to be the container's name/id if no flag are added before of the name.
func (e *Podman) Stop(args []string, attach Attach) error {
	params := []string{"stop"}
	params = append(params, args...)

	return PrintCommandR("Stopping container using the following arguments:\n  - %s", e.cmd(params, attach)).Run()
}

// Remove removes a container and returns an error in case of failure. In arguments, the first argument has to be the container's name/id if no flag are added before of the name.
func (e *Podman) Remove(args []string, attach Attach) error {
	e.Stop(args, attach)

	params := []string{"rm"}
	params = append(params, args...)

	return PrintCommandR("Removing container using the following arguments:\n  - %s", e.cmd(params, attach)).Run()
}

// Copy copies files into the container
func (e *Podman) Copy(args []string, attach Attach) *exec.Cmd {
	e.Stop(args, attach)

	params := []string{"cp"}
	params = append(params, args...)

	return PrintCommandR("Running copy using the following arguments:\n  - %s", e.cmd(params, attach))
}

// IsDocker checks if the podman path contains the word "docker". Mainly here so we can make this work on docker too.
func (e *Podman) IsDocker() bool {
	return strings.Contains(e.path, "docker")
}

// Gets the current podman version
func (e *Podman) Version() (major, minor, patch int64, err error) {
	data, err := e.cmd([]string{"--version"}, Attach{}).Output()
	if err != nil {
		return 0, 0, 0, err
	}

	regex, err := regexp.Compile(`([0-9]+)\.([0-9]+)\.([0-9]+)([0-9a-zA-z-\.]+)*`)
	if err != nil {
		return 0, 0, 0, err
	}

	parsed := regex.FindStringSubmatch(string(data))

	if parsed == nil {
		return 0, 0, 0, errors.New("unable to parse version")
	}

	major, err = strconv.ParseInt(parsed[1], 10, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	minor, err = strconv.ParseInt(parsed[2], 10, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	patch, err = strconv.ParseInt(parsed[3], 10, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

// Build builds a new image using Podman
func (e *Podman) Build(path string, tag string, attach Attach) *exec.Cmd {
	params := []string{"build"}
	params = append(params, tag)

	return PrintCommandR("Running build using the following arguments:\n  - %s", e.cmd(params, attach))
}

// Attach attaches to the podman container
func (e *Podman) Attach(args []string, attach Attach) *exec.Cmd {
	params := []string{"attach"}
	params = append(params, args...)

	return PrintCommandR("Running attach using the following arguments:\n  - %s", e.cmd(params, attach))
}

// RawCommand runs any podman subcommand (for example: ps)
func (e *Podman) RawCommand(args []string, attach Attach) *exec.Cmd {
	return e.cmd(args, attach)
}

// IsRunning checks if the container is running
func (e *Podman) IsRunning(name string) bool {
	params := []string{"container", "inspect", "-f", "'{{.State.Running}}'", name}
	data, err := e.cmd(params, Attach{}).Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(data), "true")
}
