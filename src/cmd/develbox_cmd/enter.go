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

var root bool

var enter = &cobra.Command{
	Use:   "enter",
	Short: "Enters to the container",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig()
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}
		develbox.StartContainer(configs.Podman)
		develbox.EnterContainer(&configs, root)
	},
}

func init() {
	enter.Flags().BoolVarP(&root, "root", "r", false, "Enters the container with the root user")
	rootCli.AddCommand(enter)
}
