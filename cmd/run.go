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
	Run = &cobra.Command{
		Use:   "run",
		Short: "Runs the command defined in the config file",
		Long: `Runs the command defined in the config file.
		
		Any command that is prefixed with a # inside the config will run as root.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			cfg, err := config.Read()
			if err != nil {
				return err
			}

			pman := podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})

			if _, ok := cfg.Commands[args[0]]; !ok {
				glg.Fatal("Command does not exist")
			}
			runArg := cfg.Commands[args[0]]

			// Detect if command is a string or an array of strings, and if it is prefixed with a #.
			// Here so we can know if we need to run the command as root or not.
			if _, ok := runArg.(string); ok {
				rootOpert := strings.HasPrefix(runArg.(string), "#")
				runArg = strings.TrimPrefix(runArg.(string), "#")

				params := []string{cfg.Podman.Container.Name, runArg.(string)}

				return pman.Exec(params, cfg.Image.EnvVars, true, rootOpert,
					podman.Attach{
						Stdin:     true,
						Stdout:    true,
						Stderr:    true,
						PseudoTTY: true,
					}).Run()
			}
			if _, ok := runArg.([]interface{}); ok {
				return runCommandList(pman, cfg, runArg.([]interface{}))
			}

			return glg.Errorf("\"%s\" uses an unsupported type, expected string or string array.", args[0])
		},
	}
)

func runCommandList(pman podman.Podman, cfg config.Struct, runArg []interface{}) error {
	for _, v := range runArg {
		rootOpert := strings.HasPrefix(v.(string), "#")
		newArg := strings.TrimPrefix(v.(string), "#")

		params := []string{cfg.Podman.Container.Name, newArg}

		return pman.Exec(params, cfg.Image.EnvVars, true, rootOpert,
			podman.Attach{
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				PseudoTTY: true,
			}).Run()
	}
	return nil
}
