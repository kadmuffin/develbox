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
	"os"
	"os/exec"

	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	commit = &cobra.Command{
		Use:   "commit ...",
		Short: "Commits an image.",
		Long: `Commits an image using podman commit.

Any extra arguments will get passed down to podman.
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StartContainer(configs.Podman)
			command := exec.Command(configs.Podman.Path, append([]string{"commit", configs.Podman.Container.Name}, args...)...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Run()
			develbox.StopContainer(configs.Podman)
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCli.AddCommand(commit)
}
