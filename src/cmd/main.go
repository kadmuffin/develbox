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

	"github.com/kadmuffin/develbox/src/cmd/pkg"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/spf13/cobra"
)

var (
	rootCli = &cobra.Command{
		Use:   "develbox",
		Short: "Develbox - Simple CLI tool useful for managing dev enviroments.",
		Long: `Develbox - A simple but dirty CLI tool that manages containerized dev environments.

Created so I don't have to expose my entire computer to random node modules (and to learn Go, that means BAD CODE).`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
)

func Execute() {
	// Package manager operations
	rootCli.AddCommand(pkg.Add)
	rootCli.AddCommand(pkg.Del)
	rootCli.AddCommand(pkg.Update)
	rootCli.AddCommand(pkg.Upgrade)
	rootCli.AddCommand(pkg.Search)

	if !podman.InsideContainer() {
		rootCli.AddCommand(Enter)
		rootCli.AddCommand(Create)
	}

	rootCli.Execute()
}
