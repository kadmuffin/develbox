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

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kadmuffin/develbox/cmd/create"
	"github.com/kadmuffin/develbox/cmd/dockerfile"
	"github.com/kadmuffin/develbox/cmd/pkg"
	"github.com/kadmuffin/develbox/cmd/state"
	"github.com/kadmuffin/develbox/cmd/version"
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	socketExperiment bool
	rootCli          = &cobra.Command{
		Use:     "develbox",
		Version: version.Number,
		Short:   "Develbox - CLI tool useful for creating dev environments.",
		Long: `Develbox - A CLI tool that manages containerized dev environments.

Created so I don't have to expose my entire computer to random node modules.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
)

func Execute() {
	if os.Getuid() == 0 && !podman.InsideContainer() {
		glg.Fatal("Develbox doesn't currently support being ran as root.")
	}

	// Package manager operations are added if:
	// - The user is outside the container
	// - The user is inside the container and the socket experiment is enabled
	// - The user is inside the container and running as root
	if !podman.InsideContainer() || config.FileExists(fmt.Sprintf("/home/%s/.develbox.sock", os.Getenv("USER"))) || os.Getuid() == 0 {
		// Add all the subcommands
		rootCli.AddCommand(pkg.Add)
		rootCli.AddCommand(pkg.Del)
		rootCli.AddCommand(pkg.Update)
		rootCli.AddCommand(pkg.Upgrade)
		rootCli.AddCommand(pkg.Search)
	}

	if !podman.InsideContainer() {
		socketExperiment, _ = strconv.ParseBool(os.Getenv("DEVELBOX_EXPERIMENTAL"))
		if config.ConfigExists() {
			cfg, _ := config.Read()

			if cfg.Podman.Container.Experiments.Socket {
				socketExperiment = true
			}
		}

		if socketExperiment {
			rootCli.AddCommand(Socket)
		}
		rootCli.AddCommand(Attach)
		rootCli.AddCommand(StateCmd)
		rootCli.AddCommand(Enter)

		rootCli.AddCommand(create.Create)
		rootCli.AddCommand(Exec)
		rootCli.AddCommand(Run)
		rootCli.AddCommand(state.Start)
		rootCli.AddCommand(state.Stop)
		rootCli.AddCommand(state.Restart)
		rootCli.AddCommand(state.Trash)
	}
	rootCli.AddCommand(version.VersionCmd)
	rootCli.AddCommand(dockerfile.Build)

	rootCli.Execute()
}
