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

// Package container contains the logic for creating and entering containers.
package container

import (
	"fmt"
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	globalData "github.com/kadmuffin/develbox/pkg/global"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
)

var createEtcPwd bool

// PkgVersion specifies the version of develbox (inside the container)
var PkgVersion = "latest"

// DontStopOnFinish prevents the container from stopping after the create commands are ran
var DontStopOnFinish = false

// Create creates a container and runs the setupContainer function
func Create(cfg config.Structure, deleteOld bool) error {
	pman := podman.New(cfg.Podman.Path)
	majorV, minorV, _, err := pman.Version()

	if err != nil {
		return glg.Errorf("Can't parse podman version: %s", err)
	}

	if deleteOld {
		glg.Debug("Deleting old container!")
		pman.Remove([]string{cfg.Container.Name}, podman.Attach{})
	}

	if pman.Exists(cfg.Container.Name) {
		return glg.Fail("Container already exists!")
	}

	if pman.IsDocker() {
		glg.Warn("Be aware that while probably Docker works, it may have unknown issues.")
	}

	user := os.Getenv("USER")
	uid := os.Getuid()

	args := []string{"--name", cfg.Container.Name, "-d"}

	glg.Debugf("rootless is set to: %t", cfg.Podman.Rootless)
	if cfg.Podman.Rootless {

		if !pman.IsDocker() && uid != 0 {
			// Remaps the container UID & GID so we can modify the /code folder
			args = append(args, "--userns=keep-id")
		}

		if !pman.IsDocker() && (majorV >= 5 || majorV >= 4 && minorV >= 2) {
			// Setups a user account with the current account name
			args = append(args, fmt.Sprintf("--passwd-entry=%s:*:$UID:0:develbox_container:/home/%s:/bin/sh", user, user))
		} else {
			createEtcPwd = true
		}

		// Mounts Wayland, XOrg, Pulseaudio, etc...
		xdgRunt, found := getXDGRuntime()
		if found {
			args = append(args, mountBindings(cfg, xdgRunt)...)
		} else {
			glg.Warn("Can't mount $XDG_RUNTIME_DIR, directory doesn't exist!")
		}

	}

	// Shares a global folder with the container.
	// What this means is that we can mantain certain files
	// between containers. Mainly, it's useful for
	// cache files, like nix, npm, etc...
	bindSharedFolders(cfg, &args)

	// Creates & mounts a home directory so we can access it easily
	err = os.Mkdir(".develbox/home", 0755)
	args = MountArg(args, fmt.Sprintf("%s/.develbox/home:/home/%s", config.GetCurrentDirectory(), user), false, "rslave")

	if config.GetCurrentDirectory() == os.Getenv("HOME") {
		return glg.Error("You can't create a develbox project on $HOME! (That means relabelling /home/$USER which is not a good idea)")
	}

	if err != nil && !os.IsExist(err) {
		return glg.Error("Something went wrong while creating the .develbox/home folder. %s", err)
	}

	if len(cfg.Container.Mounts) > 0 {
		args = append(args, processMounts(cfg)...)
	}
	if len(cfg.Container.Ports) > 0 {
		args = append(args, processPorts(cfg))
	}
	if len(cfg.Podman.Args) > 0 {
		args = append(args, cfg.Podman.Args...)
	}

	if cfg.Podman.Privileged {
		args = append(args, "--privileged")
	}

	args = append(args, "-e", "DEVELBOX_CONTAINER=1")
	args = append(args, getEnvVars(dfltEnvVars)...)

	if len(cfg.Container.Binds.Variables) > 0 {
		args = append(args, getEnvVars(cfg.Container.Binds.Variables)...)
	}

	// Mount configs from host
	args = MountArg(args, "/etc/localtime:/etc/localtime", true, "")
	args = MountArg(args, "/etc/resolv.conf:/etc/resolv.conf", true, "")
	args = MountArg(args, "/etc/hosts:/etc/hosts", true, "")
	args = MountArg(args, "/etc/timezone:/etc/timezone", true, "")
	args = MountArg(args, "/home/%s/.gitconfig:/etc/gitconfig", true, "")

	// Add DEVELBOX_VERSION label to the container
	args = append(args, "--label", fmt.Sprintf("develbox_version=%s", PkgVersion))
	// Add DEVELBOX_CONTAINER label to the container
	args = append(args, "--label", "develbox_container=1")
	// Add DEVELBOX_PROJECT_PATH label to the container
	args = append(args, "--label", fmt.Sprintf("develbox_project_path=\"%s\"", config.GetCurrentDirectory()))

	// Mount the main folder and pass the image URI before the
	// container is created.
	args = append(args, mountWorkDir(cfg)...)
	args = append(args, cfg.Image.URI, "sh")

	err = pman.Create(args, podman.Attach{Stdout: true, Stderr: true}).Run()

	if !pman.IsRunning(cfg.Container.Name) {
		glg.Warnf("Container '%s' is not running!.", cfg.Container.Name)
	}

	if err != nil {
		return glg.Error("Something went wrong while creating the container.")
	}
	// Adds the current user to /etc/passwd
	// Only used if the current podman version doesn't
	// support --passwd-entry
	if createEtcPwd {
		cmd := pman.Exec([]string{cfg.Container.Name, fmt.Sprintf("echo '%s:*:%d:0:develbox_container:/home/%s:/bin/sh' >> /etc/passwd", user, os.Getuid(), user)}, map[string]string{}, true, true, podman.Attach{Stderr: true})
		glg.Debugf("Running command: %s", cmd.String())
		err = cmd.Run()

		if err != nil {
			glg.Warnf("Couldn't setup user on /etc/passwd: %s", err)
		}
	}

	setupContainer(&pman, cfg)

	if !DontStopOnFinish {
		pman.Stop([]string{cfg.Container.Name}, podman.Attach{Stderr: true})
	}
	DontStopOnFinish = false

	if cfg.Podman.AutoCommit {
		glg.Warn("Auto commit feature is enabled, deleting old image (if exists) and commiting new one.")

		// Deletes the old image
		pman.RawCommand([]string{"rmi", cfg.Container.Name}, podman.Attach{Stderr: true})

		// Commit new image
		pman.Commit([]string{cfg.Container.Name, cfg.Image.URI}, podman.Attach{Stderr: true})
	}
	if cfg.Podman.AutoDelete {
		glg.Warn("Auto delete feature is enabled, deleting container.")
		pman.Remove([]string{cfg.Container.Name}, podman.Attach{Stderr: true})
	} else {
		fmt.Println("Enter to the container with: develbox enter.")
	}

	fmt.Println("Operation completed!")
	return nil
}

// setupContainer installs the packages and runs the onCreation & onFinish commands
func setupContainer(pman *podman.Podman, cfg config.Structure) {
	// Runs commands that should be ran just
	// after the container was created
	err := RunCommandList(cfg.Container.Name,
		cfg.Image.OnCreation, pman, true,
		podman.Attach{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			PseudoTTY: true,
		})

	if err != nil {
		pman.Remove([]string{cfg.Container.Name}, podman.Attach{Stderr: true})
		glg.Fatalf("Something went wrong with creating your container. %s", err)
	}

	opert := pkgm.NewOperation("update", []string{}, []string{}, true)
	opert.DevInstall = true
	opert.Process(&cfg)

	var goInstalled bool

	// Check if cfg.Packages or cfg.DevPackages contain go
	if !Contains(cfg.Packages, "go") && !Contains(cfg.DevPackages, "go") {
		fmt.Println("> Installing Go for develbox experimental features")
		opert := pkgm.NewOperation("add", []string{"go"}, []string{}, true)
		cmd, _ := opert.ProcessCmd(&cfg, podman.Attach{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		})

		if cmd.Run() != nil {
			glg.Warn("Couldn't install Go, some features may not work")
		} else {
			goInstalled = true
		}
	}

	if len(cfg.Packages)+len(cfg.DevPackages) > 0 {
		installPkgs(pman, cfg, append(cfg.Packages, cfg.DevPackages...), true)
	}

	if len(cfg.UserPkgs.Packages)+len(cfg.UserPkgs.DevPackages) > 0 && cfg.Podman.Rootless {
		err := installPkgs(pman, cfg, append(cfg.UserPkgs.Packages, cfg.UserPkgs.DevPackages...), false)

		if err != nil {
			pman.Remove([]string{cfg.Container.Name}, podman.Attach{Stderr: true})
			glg.Fatalf("Something went wrong while installing the specified packages. %s", err)
		}
	}

	if goInstalled {
		fmt.Println("> Installing develbox inside the container")
		err := RunCommandList(cfg.Container.Name, []string{fmt.Sprintf("go install github.com/kadmuffin/develbox@v%s", PkgVersion), "cp /root/go/bin/develbox /usr/local/bin/develbox"}, pman, true, podman.Attach{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		})

		if err != nil {
			glg.Warnf("An error occurred and develbox wasn't installed. %s", err)
		}
	}

	err = RunCommandList(cfg.Container.Name,
		cfg.Image.OnFinish, pman, true,
		podman.Attach{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			PseudoTTY: true,
		})

	if err != nil {
		pman.Remove([]string{cfg.Container.Name}, podman.Attach{Stderr: true})
		glg.Fatalf("Something went wrong with finishing setting up your container. %s", err)
	}
}

func installPkgs(pman *podman.Podman, cfg config.Structure, pkgs []string, root bool) error {
	opert := pkgm.NewOperation("add", pkgs, []string{}, true)
	opert.UserOperation = !root
	cmd, _ := opert.ProcessCmd(&cfg, podman.Attach{Stdin: true, Stdout: true, Stderr: true})
	return cmd.Run()
}

// DontAttachEnter is a flag that tells the enter command to not attach to the container (for the enter command)
var DontAttachEnter = false

// Enter runs a shell in the container and creates a pipe for package installations.
func Enter(cfg config.Structure, root bool) error {
	pman := podman.New(cfg.Podman.Path)

	attach := podman.Attach{
		Stdin:     !DontAttachEnter,
		Stdout:    !DontAttachEnter,
		Stderr:    !DontAttachEnter,
		PseudoTTY: !DontAttachEnter,
	}

	cmd := pman.Exec([]string{cfg.Container.Name, cfg.Container.Shell}, cfg.Image.Variables, false, root, attach)

	if DontAttachEnter {
		DontAttachEnter = false
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", err, out)
		}
		return nil
	}
	return cmd.Run()
}

// InstallAndEnter install the packages and runs a shell in the container
func InstallAndEnter(cfg config.Structure, root bool) error {
	pman := podman.New(cfg.Podman.Path)

	err := installPkgs(&pman, cfg, append(cfg.Packages, cfg.DevPackages...), true)
	if err != nil {
		return glg.Errorf("Couldn't install packages. %s", err)
	}

	return Enter(cfg, root)
}

// Loops through the shared folders and creates and binds them to the container.
func bindSharedFolders(cfg config.Structure, args *[]string) {
	for key, value := range cfg.Container.SharedFolders {
		switch value := value.(type) {
		case string:
			endPath := ReplaceEnvVars(value)
			newPath, err := globalData.CreateFile(endPath, key)

			if err != nil {
				glg.Fatalf("Couldn't create the shared folder %s. %s", endPath, err)
			}
			*args = append(*args, fmt.Sprintf("-v=%s:%s:rw,z", newPath, endPath))
		case []interface{}:
			for _, val := range value {
				endPath := ReplaceEnvVars(val.(string))
				newPath, err := globalData.CreateFile(endPath, key)

				if err != nil {
					glg.Fatalf("Couldn't create the shared folder %s. %s", endPath, err)
				}

				*args = append(*args, fmt.Sprintf("-v=%s:%s:rw,z", newPath, endPath))
			}
		default:
			glg.Fatalf("The shared folder value must be a string or a list of strings. Got %d", value)
		}
	}
}
