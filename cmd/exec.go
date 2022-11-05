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

package cmd

import (
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	Exec = &cobra.Command{
		Use:                "exec",
		Short:              "Executes a program inside the container",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cfg, err := config.Read()
			if err != nil {
				return err
			}
			pman := podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})

			var rootOpert bool
			joinedArgs := strings.Join(args, " ")
			if strings.HasPrefix(joinedArgs, "#") {
				rootOpert = true
				joinedArgs = strings.TrimPrefix(joinedArgs, "#")
			}
			if strings.HasPrefix(joinedArgs, "!") {
				rootOpert = true
				joinedArgs = strings.TrimPrefix(joinedArgs, "!")
			}

			params := []string{cfg.Podman.Container.Name, joinedArgs}
			command := pman.Exec(params, cfg.Image.EnvVars, true, rootOpert,
				podman.Attach{
					Stdin:     true,
					Stdout:    true,
					Stderr:    true,
					PseudoTTY: true,
				})

			return command.Run()
		},
	}
)
