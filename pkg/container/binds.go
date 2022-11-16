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

package container

import (
	"fmt"
	"os"
	"strings"

	"os/exec"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
)

var dfltEnvVars = []string{
	"XDG_CURRENT_DESKTOP",
	"XDG_RUNTIME_DIR",
	"XDG_SESSION_CLASS",
	"XDG_SESSION_DESKTOP",
	"XDG_SESSION_TYPE",
	"DBUS_SESSION_BUS_ADDRESS",
	"DESKTOP_SESSION",
	"WAYLAND_DISPLAY",
	"DISPLAY",
	"GDK_BACKEND",
	"PULSE_SERVER",
	"USER",
}

// Returns a string list with the enviroment variables to copy into the container.
func getEnvVars(vars []string) []string {
	// See https://github.com/containers/toolbox/blob/main/src/pkg/utils/utils.go#L273
	// for where this comes from.
	result := []string{}
	for _, envVar := range vars {
		variable, found := os.LookupEnv(envVar)
		if found {
			result = append(result, "-e", fmt.Sprintf("%s=%s", envVar, variable))
			continue
		}
		glg.Debugf("Didn't find enviroment var '%s'.", envVar)
	}

	return result
}

// Mounts the directory used by xorg and adds user with xhost
func mountXOrg(args []string) []string {
	if !config.FileExists("/tmp/.X11-unix") {
		glg.Errorf("didn't find XOrg socket, skipping mount...")
		return args
	}

	glg.Debugf("Mounting XOrg.")
	exec.Command("xhost", fmt.Sprintf("+SI:localhost:%s", os.Getenv("USER"))).Run()
	return append(args, "-v=/tmp/.X11-unix:/tmp/.X11-unix:rslave")
}

// Mounts /dev with the rslave option.
func mountDev(cfg config.Podman) []string {
	// We mount /dev so we can access things like cameras and GPUs
	// inside the container. See github.com/containers/podman/issues/5623.
	if strings.Contains(cfg.Path, "docker") {
		return []string{"-v=/dev:/dev:rslave"}
	}

	return []string{
		"-v=/dev/:/dev:rslave",
		"--mount", "type=devpts,destination=/dev/pts",
	}

}

// Mounts the Workspace directory with proper SELinux label if necessary.
func mountWorkDir(cfg config.Structure) []string {
	workDir := cfg.Container.WorkDir
	mntOpts := ""

	mountString := []string{"--mount", fmt.Sprintf("type=bind,source=%s,destination=%s,bind-propagation=rslave", config.GetCurrentDirectory(), workDir)}

	// Adding "private unshare label" so SELinux doesn't
	// get mad at us when running without "--privileged"
	if cfg.Podman.Rootless && !cfg.Podman.Privileged {
		mntOpts += ":Z"

		mountString = []string{fmt.Sprintf("-v=%s:%s%s", config.GetCurrentDirectory(), workDir, mntOpts)}
	}
	return append(mountString, fmt.Sprintf("-w=%s", workDir))

}

// Returns a string pointing to $XDG_RUNTIME_DIR and a boolean indicating if the folder exists or not.
func getXDGRuntime() (string, bool) {
	xdgRuntime, found := os.LookupEnv("XDG_RUNTIME_DIR")

	if !found {
		xdgRuntime = fmt.Sprintf("/run/user/%d", os.Getuid())
		glg.Warnf("Couldn't find $XDG_RUNTIME_DIR, assuming %s", xdgRuntime)
	}

	return xdgRuntime, config.FileExists(xdgRuntime)
}

// Mounts all the required binds in the config file.
func mountBindings(cfg config.Structure, xdgRuntime string) []string {
	args := []string{}

	os.Setenv("XDG_RUNTIME_DIR", xdgRuntime)

	if cfg.Container.Binds.XOrg {
		args = mountXOrg(args)
	}

	if cfg.Container.Binds.Dev {
		args = append(args, mountDev(cfg.Podman)...)
	}

	args = append(args, fmt.Sprintf("-v=%s:%s:rslave", xdgRuntime, xdgRuntime))

	return args
}
