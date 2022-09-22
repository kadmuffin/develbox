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

	"github.com/creasty/defaults"
	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var writeConfig bool

var create = &cobra.Command{
	Use:   "create",
	Short: "Creates a container using the develbox.json",
	Run: func(cmd *cobra.Command, args []string) {

		if writeConfig {
			configs := new(develbox.DevSetings)
			develbox.SetContainerName(configs)
			defaults.Set(configs)
			os.Mkdir(".develbox", 0755)
			if develbox.ConfigExists() && !forceAction {
				log.Fatal("A config file already exist!")
			}
			os.Remove(".develbox/config.json")
			develbox.WriteConfig(configs)
			os.Exit(0)
		}
		var configs develbox.DevSetings = develbox.ReadConfig()
		develbox.CreateContainer(&configs, forceAction)
	},
}

func init() {
	create.Flags().BoolVarP(&forceAction, "force", "f", false, "Forces to overwrite the container/config")
	create.Flags().BoolVarP(&writeConfig, "config", "c", false, "Writes a config file with pre-set defaults")
	rootCli.AddCommand(create)
}
