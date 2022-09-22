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
	"github.com/spf13/cobra"
)

var (
	forceAction bool
	rootCli     = &cobra.Command{
		Use:   "develbox",
		Short: "Develbox - Simple CLI tool useful for managing dev enviroments.",
		Long: `Develbox - A simple but dirty CLI tool that manages containerized dev enviroments.

Created so I don't have to expose my entire computer to random node modules (and to learn Go, that means BAD CODE).`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Execute() {
	rootCli.PersistentFlags().BoolVarP(&forceAction, "force", "f", false, "Forces the subsequent action to execute.")
	rootCli.Execute()
}
