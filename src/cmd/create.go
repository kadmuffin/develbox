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
	"github.com/spf13/cobra"
)

var (
	createCfg    bool
	forceReplace bool
	Create       = &cobra.Command{
		Use:        "create",
		SuggestFor: []string{"config"},
		Short:      "Creates a new container/config for this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createCfg {
				cfg := config.Struct{}
				cfg.SetDefaults()
				return config.WriteConfig(&cfg)
			}

			cfg, err := config.Read()
			if err != nil {
				return err
			}
			container.Create(cfg)
			return nil
		},
	}
)

func init() {
	Enter.Flags().BoolVarP(&createCfg, "config", "c", false, "Use to create a new config file")
	Enter.Flags().BoolVarP(&forceReplace, "force", "f", false, "Use to force the creation of a container/config")
}
