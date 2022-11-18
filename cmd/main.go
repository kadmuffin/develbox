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

// Package cmd contains the some commands for the program
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
	// rootCLI is the root command for the program
	rootCLI = &cobra.Command{
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

// Execute is the entrypoint for the program
func Execute() error {
	if os.Getuid() == 0 && !podman.InsideContainer() {
		glg.Fatal("Develbox doesn't currently support being ran as root.")
	}

	// Package manager operations are added if:
	// - The user is outside the container
	// - The user is inside the container and the socket experiment is enabled
	// - The user is inside the container and running as root
	if !podman.InsideContainer() || config.FileExists(fmt.Sprintf("/home/%s/.develbox.sock", os.Getenv("USER"))) || os.Getuid() == 0 {
		// Add all the subcommands
		rootCLI.AddCommand(pkg.Add)
		rootCLI.AddCommand(pkg.Del)
		rootCLI.AddCommand(pkg.Update)
		rootCLI.AddCommand(pkg.Upgrade)
		rootCLI.AddCommand(pkg.Search)
	}

	if !podman.InsideContainer() {
		socketExperiment, _ = strconv.ParseBool(os.Getenv("DEVELBOX_EXPERIMENTAL"))
		if config.Exists() {
			cfg, _ := config.Read()

			if cfg.Experiments.Socket {
				socketExperiment = true
			}
		}

		if socketExperiment {
			rootCLI.AddCommand(Socket)
		}
		rootCLI.AddCommand(Attach)
		rootCLI.AddCommand(state.Cmd)
		rootCLI.AddCommand(Enter)

		rootCLI.AddCommand(create.Create)
		rootCLI.AddCommand(Exec)
		rootCLI.AddCommand(Run)
		rootCLI.AddCommand(state.Start)
		rootCLI.AddCommand(state.Stop)
		rootCLI.AddCommand(state.Restart)
		rootCLI.AddCommand(state.Trash)
	}
	rootCLI.AddCommand(version.VersionCmd)
	rootCLI.AddCommand(dockerfile.Build)

	return rootCLI.Execute()
}

// GetRootCLI returns the root command for the program
func GetRootCLI() *cobra.Command {
	return rootCLI
}
