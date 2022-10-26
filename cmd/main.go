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
	"os"

	"github.com/kadmuffin/develbox/cmd/create"
	"github.com/kadmuffin/develbox/cmd/dockerfile"
	"github.com/kadmuffin/develbox/cmd/pkg"
	"github.com/kadmuffin/develbox/cmd/state"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	rootCli = &cobra.Command{
		Use:     "develbox",
		Version: "0.1.2",
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
	if os.Getuid() == 0 {
		glg.Fatal("Develbox doesn't currently support being ran as root.")
	}

	if !podman.InsideContainer() {
		// Package manager operations
		rootCli.AddCommand(pkg.Add)
		rootCli.AddCommand(pkg.Del)
		rootCli.AddCommand(pkg.Update)
		rootCli.AddCommand(pkg.Upgrade)
		rootCli.AddCommand(pkg.Search)

		rootCli.AddCommand(Enter)
		rootCli.AddCommand(create.Create)
		rootCli.AddCommand(Exec)
		rootCli.AddCommand(Run)
		rootCli.AddCommand(state.Start)
		rootCli.AddCommand(state.Stop)
		rootCli.AddCommand(state.Trash)
	}
	rootCli.AddCommand(Version)
	rootCli.AddCommand(dockerfile.Build)

	rootCli.Execute()
}
