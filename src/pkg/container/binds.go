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

	"os/exec"

	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/kpango/glg"
)

var envVars = []string{
	"XDG_CURRENT_DESKTOP",
	"XDG_RUNTIME_DIR",
	"XDG_SESSION_CLASS",
	"XDG_SESSION_DESKTOP",
	"XDG_SESSION_TYPE",
	"DBUS_SESSION_BUS_ADDRESS",
	"DESKTOP_SESSION",
	"WAYLAND_DISPLAY",
	"DISPLAY",
	// Not adding XAuthority as I don't know a way other than not mount the
	// entire /run/user/$UID directory in the container
	//	"XAUTHORITY",
	"GDK_BACKEND",
	"PULSE_SERVER",
	"USER",
}

// Returns a string list with the enviroment variables to copy
// into the container.
//
// See https://github.com/containers/toolbox/blob/main/src/pkg/utils/utils.go#L273
func getEnvVars() []string {
	result := []string{}
	for _, envVar := range envVars {
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
func mountXOrg() string {
	if !config.FileExists("/tmp/.X11-unix") {
		glg.Errorf("didn't find xorg socket, skipping mount...")
		return ""
	}

	glg.Debugf("Mounting XOrg.")
	exec.Command("xhost", fmt.Sprintf("+SI:localhost:%s", os.Getenv("USER"))).Run()
	return "-v=/tmp/.X11-unix:/tmp/.X11-unix:rw"
}

// Loads the Wayland socket into the container
func mountWayland(xdrRunt string) string {
	waylandDisplay, found := os.LookupEnv("WAYLAND_DISPLAY")

	if !found {
		glg.Errorf("didn't find wayland display, skipping mount...")
		return ""
	}
	glg.Debugf("Mounting Wayland socket: %s/%s", xdrRunt, waylandDisplay)

	return fmt.Sprintf("-v=%s/%s:%s/%s", xdrRunt, waylandDisplay, xdrRunt, waylandDisplay)
}

// Mounts the pipewire socket
//
// Here so we can have audio inside the container
// using pipewire directly.
func mountPipepwire(xdrRunt string) string {
	pwPath := fmt.Sprintf("%s/pipewire-0", xdrRunt)

	if !config.FileExists(pwPath) {
		glg.Errorf("didn't find pipewire socket, skipping mount...")
		return ""
	}

	glg.Debugf("Mounting Pipewire socket: %s", pwPath)
	return fmt.Sprintf("-v=%s:%s", pwPath, pwPath)
}

// Mounts the pulseaudio socket
//
// Here so we can have audio inside the container.
// As a lot of distros don't use pipewire, so we load pulseaudio
// by default. Pipewire offers a pulseaudio socket so it
// should also work.
func mountPulseaudio(xdrRunt string) []string {
	paPath := fmt.Sprintf("%s/pulse/native", xdrRunt)
	if !config.FileExists(paPath) {
		glg.Errorf("didn't find the pulseaudio socket, skipping mount...")
		return []string{}
	}
	return []string{
		fmt.Sprintf("-v=%s:%s", paPath, paPath),
		"--device=/dev/snd",
	}
}

// Mounts /dev with the rslave option.
//
// We mount /dev so we can access things like cameras and GPUs
// inside the container. See github.com/containers/podman/issues/5623.
func mountDev() []string {
	return []string{
		"-v=/dev:/dev/rslave",
		"--mount", "type=devpts,dest=/dev/pts",
	}
}

// Mounts the Workspace directory with proper SELinux label
// if necessary.
func mountWorkDir(cfg config.Struct) string {
	workDir := cfg.Podman.Container.WorkDir

	// Adding "private unshare label" so SELinux doesn't
	// get mad at us when running without "--privileged"
	if cfg.Podman.Rootless {
		workDir += ":Z"
	}

	return fmt.Sprintf("-v=%s:%s -w=%s", config.GetCurrentDirectory(), workDir, workDir)
}

// Copy develbox from $GOPATH/bin
func copyDevBBin(pman podman.Podman, cname string) {
	value, found := os.LookupEnv("GOPATH")
	if !found {
		home := os.Getenv("HOME")
		value = fmt.Sprintf("%s/go", home)
	}

	pman.Copy([]string{fmt.Sprintf("%s/bin/develbox", value), fmt.Sprintf("%s:/bin/develbox", cname)}, podman.Attach{})
}

// Mounts all the required binds in the config file.
//
// Wayland & Pulseaudio get mounted by default. And possible
// bindings are XOrg, Pipewire, /dev (so we can access cameras & DRI).
func mountBindings(cfg config.Struct) []string {
	args := []string{}

	xdrRuntime, found := os.LookupEnv("XDR_RUNTIME_DIR")

	if !found {
		xdrRuntime = fmt.Sprintf("/run/user/%d", os.Getuid())
		os.Setenv("XDG_RUNTIME_DIR", xdrRuntime)
		glg.Warnf("Couldn't find $XDR_RUNTIME_DIR, assuming %s", xdrRuntime)
	}

	if cfg.Podman.Container.Binds.XOrg {
		args = append(args, mountXOrg())
	}

	args = append(args, mountWayland(xdrRuntime))
	args = append(args, mountPulseaudio(xdrRuntime)...)

	if cfg.Podman.Container.Binds.Pipewire {
		args = append(args, mountPipepwire(xdrRuntime))
	}

	if cfg.Podman.Container.Binds.Dev {
		args = append(args, mountDev()...)
	}

	return args
}
