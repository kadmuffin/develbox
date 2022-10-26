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

package pkg

import (
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	Del = &cobra.Command{
		Use:                "del",
		SuggestFor:         []string{"delete", "remove"},
		Short:              "Deletes packages from the container",
		Long:               "Deletes packages using the package manager defined in the config.",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			packages, flags := pkgm.ParseArguments(args)
			parsedFlags := parseFlags(flags)

			if parsedFlags.ShowHelp {
				cmd.Help()
				return
			}

			opertn := pkgm.NewOperation("del", *packages, parsedFlags.All, false)
			opertn.UserOperation = parsedFlags.UserOpert

			cfg, err := config.Read()
			if err != nil {
				glg.Error(err)
				return
			}
			pman := podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}

			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})
			opertn.Process(&cfg, false)
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
	Del.Flags().BoolP("user", "U", false, "Install packages as user instead of root.")
	Del.Flags().BoolP("pkg-help", "p", false, "Show the package manager help for this command.")

}
