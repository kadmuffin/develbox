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
	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/pkgm"
	"github.com/kadmuffin/develbox/src/pkg/podman"
	"github.com/spf13/cobra"
)

var (
	Search = &cobra.Command{
		Use:                "search",
		Aliases:            []string{"srch"},
		Short:              "Search for packages using the pkg manager",
		Long:               "Search for (all, usually) matching packages using the package manager defined in the config.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			packages, flags := pkgm.ParseArguments(args)
			opertn := pkgm.NewOperation("search", *packages, *flags)

			cfg, err := config.Read()
			if err != nil {
				return err
			}
			pman := podman.New(cfg.Podman.Path)
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})
			return opertn.Process(&cfg)
		},
	}
)
