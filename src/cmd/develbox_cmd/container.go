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
	container = &cobra.Command{
		Use:   "container",
		Short: "Manage the container state with the sub-commands",
	}

	start = &cobra.Command{
		Use:   "start",
		Short: "Start the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StartContainer(configs.Podman)
		},
	}

	stop = &cobra.Command{
		Use:   "stop",
		Short: "Stop the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StopContainer(configs.Podman)
		},
	}

	delete = &cobra.Command{
		Use:   "delete",
		Short: "Delete the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StopContainer(configs.Podman)
		},
	}
)

func init() {
	container.AddCommand(start)
	container.AddCommand(stop)
	container.AddCommand(delete)

	rootCli.AddCommand(container)
}
