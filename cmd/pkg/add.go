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
	"github.com/spf13/cobra"
)

var (
	Add = &cobra.Command{
		Use:                "add",
		SuggestFor:         []string{"install"},
		Short:              "Installs packages into the container",
		Long:               "Installs packages using the package manager defined in the config.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			packages, flags := pkgm.ParseArguments(args)
			parsedFlags, devInstall := parseFlags(flags)
			opertn := pkgm.NewOperation("add", *packages, *parsedFlags, false)

			cfg, err := config.Read()
			if err != nil {
				return err
			}
			pman := podman.New(cfg.Podman.Path)
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})
			err = opertn.Process(&cfg, devInstall)
			if err != nil {
				return err
			}
			return config.Write(&cfg)
		},
	}
)

func parseFlags(flags *[]string) (*[]string, bool) {
	var devInstall bool
	var parsedFlags []string
	for _, flag := range *flags {
		if flag == "--dev" || flag == "-D" {
			devInstall = true
			continue
		}
		parsedFlags = append(parsedFlags, flag)
	}
	return &parsedFlags, devInstall
}

func init() {
	Add.Flags().BoolP("dev", "D", false, "Install packages as development dependencies")
}
