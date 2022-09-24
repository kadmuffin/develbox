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
	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/container"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	createCfg    bool
	forceReplace bool
	Create       = &cobra.Command{
		Use:        "create",
		SuggestFor: []string{"config"},
		Short:      "Creates a new container/config for this project",
		Run: func(cmd *cobra.Command, args []string) {
			if createCfg {
				cfg := config.Struct{}
				config.SetDefaults(&cfg)
				err := config.WriteConfig(&cfg)
				if err != nil {
					glg.Error(err)
				}
				return
			}

			cfg, err := config.Read()
			if err != nil {
				glg.Error(err)
				return
			}
			container.Create(cfg)
			return
		},
	}
)

func init() {
	Create.Flags().BoolVarP(&createCfg, "config", "c", false, "Use to create a new config file")
	Create.Flags().BoolVarP(&forceReplace, "force", "f", false, "Use to force the creation of a container/config")
}