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

	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/kpango/glg"
)

func Create(cfg config.Struct, deleteOld bool) {
	pman := podman.New(cfg.Podman.Path)

	if deleteOld {
		glg.Debug("Deleting old container!")
		pman.Remove([]string{cfg.Podman.Container.Name}, podman.Attach{})
	}

	if pman.Exists(cfg.Podman.Container.Name) {
		glg.Fail("Container already exists!")
		os.Exit(1)
	}

	user := os.Getenv("USER")

	args := []string{"--name", cfg.Podman.Container.Name}

	glg.Debugf("rootless is set to: %t", cfg.Podman.Rootless)
	if cfg.Podman.Rootless {

		// Because I can't figure it right now, rootless with docker won't work
		if pman.IsDocker() {
			glg.Error("Rootless mode doesn't currently work with Docker")
			os.Exit(1)
		}

		// Remaps the container UID & GID so we can modify the /code folder
		args = append(args, "--userns=keep-id")

		// Setups a user account with the current account name
		args = append(args, fmt.Sprintf("--passwd-entry=%s:*:$UID:0:develbox_container:/home/%s:/bin/sh", user, user))

		// Mounts Wayland, XOrg, Pulseaudio, etc...
		args = append(args, mountBindings(cfg)...)
	}

	// Creates & mounts a home directory so we can access it easily
	// Not usable as a home directory without rootless mode
	// We just mount it right now to get the pipe working
	err := os.Mkdir(".develbox/home", 0755)
	args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s/.develbox/home,destination=/home/%s,bind-propagation=rslave", config.GetCurrentDirectory(), user))

	if err != nil && !os.IsExist(err) {
		glg.Fatalf("Something went wrong while creating the .develbox/home folder. %s", err)
	}

	if len(cfg.Podman.Container.Mounts) > 0 {
		args = append(args, processVolumes(cfg))
	}
	if len(cfg.Podman.Container.Ports) > 0 {
		args = append(args, processPorts(cfg))
	}
	if len(cfg.Podman.Container.Args) > 0 {
		args = append(args, cfg.Podman.Container.Args...)
	}

	if cfg.Podman.Container.Privileged {
		args = append(args, "--privileged")
	}

	args = append(args, getEnvVars()...)

	// Mount the main folder and pass the image URI before the
	// container is created.
	args = append(args, mountWorkDir(cfg)...)
	args = append(args, cfg.Image.URI)
	err = pman.Create(args, podman.Attach{Stdout: true, Stderr: true}).Run()

	if err != nil {
		glg.Fatalf("Something went wrong while creating the container.")
	}

	setupContainer(&pman, cfg)

	pman.Stop([]string{cfg.Podman.Container.Name}, podman.Attach{Stderr: true})

	if cfg.Podman.BuildOnly {
		pman.Remove([]string{cfg.Podman.Container.Name}, podman.Attach{Stderr: true})
	}

	fmt.Println("Operation completed!")
	fmt.Println("Enter to the container with: develbox enter.")
}

// Sets the packages up and runs the onCreation & onFinish commands
func setupContainer(pman *podman.Podman, cfg config.Struct) {
	// Runs commands that should be ran just
	// after the container was created
	err := RunCommandList(cfg.Podman.Container.Name,
		cfg.Image.OnCreation, pman, true,
		podman.Attach{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			PseudoTTY: true,
		})

	if err != nil {
		pman.Remove([]string{cfg.Podman.Container.Name}, podman.Attach{Stderr: true})
		glg.Fatalf("Something went wrong with creating your container. %s", err)
	}

	if len(cfg.Packages) > 0 {
		spacedPkgs := strings.Join(cfg.Packages, " ")
		addBase := strings.ReplaceAll(cfg.Image.Installer.Add, " {-y}", " -y")

		command := []string{cfg.Podman.Container.Name,
			strings.ReplaceAll(addBase, "{args}", spacedPkgs)}

		err = pman.Exec(command, true, true,
			podman.Attach{
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				PseudoTTY: true,
			}).Run()

		if err != nil {
			pman.Remove([]string{cfg.Podman.Container.Name}, podman.Attach{Stderr: true})
			glg.Fatalf("Something went wrong while installing the specified packages. %s", err)
		}
	}

	err = RunCommandList(cfg.Podman.Container.Name,
		cfg.Image.OnFinish, pman, true,
		podman.Attach{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			PseudoTTY: true,
		})

	if err != nil {
		pman.Remove([]string{cfg.Podman.Container.Name}, podman.Attach{Stderr: true})
		glg.Fatalf("Something went wrong with finishing setting up your container. %s", err)
	}
}

// Runs a shell in the container
// and creates a pipe for package
// installations.
func Enter(cfg config.Struct, root bool) error {
	pman := podman.New(cfg.Podman.Path)
	//	pipe := pipes.New(".develbox/home/.develbox")
	//	pipe.Create()

	cmd := pman.Exec([]string{cfg.Podman.Container.Name, cfg.Podman.Container.Shell}, false, root,
		podman.Attach{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			PseudoTTY: true,
		})

	//	go pkgPipe(&cfg, pipe)
	err := cmd.Run()
	//	pipe.Remove()
	return err
}
