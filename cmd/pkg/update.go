// Copyright 2022 Kevin Ledesma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pkg has the logic for communication with the package manager
package pkg

import (
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	// Update is the command to update the package manager
	Update = &cobra.Command{
		Use:     "update",
		Aliases: []string{"upd"},
		Short:   "Updates the package databases in container",
		Long: `Updates package databases using the package manager defined in the config.

		To actually update packages try using upgrade instead.
		`,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			packages, flags := pkgm.ParseArguments(args)
			parsedFlags := parseFlags(&flags)

			if parsedFlags.ShowHelp || len(packages)+len(parsedFlags.All) == 0 {
				cmd.Help()
				return
			}

			opertn := pkgm.NewOperation("update", packages, parsedFlags.All, false)
			opertn.UserOperation = parsedFlags.UserOpert

			cfg, err := config.Read()
			if err != nil {
				glg.Error(err)
				return
			}
			StartContainer(&cfg)

			if podman.InsideContainer() && os.Getuid() != 0 {
				SendOperation(opertn)
				return
			}

			opertn.Process(&cfg)
			if err != nil {
				glg.Error(err)
				return
			}
			err = config.Write(&cfg)
			if err != nil {
				glg.Error(err)
			}
		},
	}
)

func init() {
	Update.Flags().BoolP("user", "U", false, "Install packages as user instead of root.")
	Update.Flags().BoolP("pkg-help", "p", false, "Show the package manager help for this command.")

}
