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

package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	rootOperation bool
	customCommand bool

	run = &cobra.Command{
		Use:   "run ...",
		Short: "Run the specified command defined in the config file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			arguments := args

			if arguments[0] == "-r" {
				rootOperation = true
				arguments = arguments[1:]
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands(configs.Commands[arguments[0]], configs.Podman, true, false, false, true, rootOperation)
		},
	}

	runc = &cobra.Command{
		Use:                "runc ...",
		Short:              "Run the custom command passed in the argument",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			arguments := args

			if arguments[0] == "-r" {
				rootOperation = true
				arguments = arguments[1:]
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommand(arguments, configs.Podman, true, false, false, true, "%s", rootOperation)
		},
	}
)

func init() {
	rootCli.AddCommand(run)
	rootCli.AddCommand(runc)
}
