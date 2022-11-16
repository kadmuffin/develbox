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

// Package state manages the state of the container
package state

import (
	"fmt"
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

// Cmd is the cobra command for the state command
var Cmd = &cobra.Command{
	Use:   "state",
	Short: "Prints the current container state",
	Long:  `Prints the current container state`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Read()

		if err != nil {
			glg.Fatalf("Can't read config file: %s", err)
		}

		pman := podman.New(cfg.Podman.Path)

		if !pman.Exists(cfg.Container.Name) {
			fmt.Println("Container does not exist.")
			os.Exit(1)
		}

		running := pman.IsRunning(cfg.Container.Name)
		if running {
			fmt.Println("Container is running!")
			os.Exit(0)
		} else {
			fmt.Println("Container exists but is not running!")
			os.Exit(2)
		}
	},
}
