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
	"fmt"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	cfg  config.Struct
	pman podman.Podman
	Run  = &cobra.Command{
		Use:   "run",
		Short: "Runs the command defined in the config file",
		Long: `Runs the command defined in the config file.
		
		Any command that is prefixed with a # inside the config will run as root. Call other commands using the "!" prefix.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			var err error
			cfg, err = config.Read()
			if err != nil {
				return err
			}

			pman = podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})

			name := strings.Join(args, " ")
			if _, ok := cfg.Commands[name]; !ok {
				return glg.Errorf("Command '%s' does not exist")
			}

			runArgs, err := getAllAsArray(name, "")
			if err != nil {
				return err
			}

			return runCommandList(runArgs)
		},
	}
)

// getAllAsArray takes a name and returns an array of strings.
//
// The resulting array will contain all the commands inside that name.
//
// If the name is prefixed with "!", it will recursively call itself
// to get the full command tree.
func getAllAsArray(name string, from string) ([]string, error) {
	if _, ok := cfg.Commands[name]; !ok {
		return []string{}, glg.Errorf("[%s] Command '%s' does not exist", from, name)
	}
	cmds := cfg.Commands[name]
	result := []string{}

	if _, ok := cmds.(string); ok {
		if strings.HasPrefix(cmds.(string), "!") {
			return parseRecursion(cmds.(string), name, from)
		}
		result = append(result, cmds.(string))
		return result, nil
	}
	if _, ok := cmds.([]interface{}); ok {
		for _, v := range cmds.([]interface{}) {
			if strings.HasPrefix(v.(string), "!") {
				newCmds, err := parseRecursion(v.(string), name, from)

				if err != nil {
					return []string{}, err
				}

				result = append(result, newCmds...)
				continue
			}
			result = append(result, v.(string))
		}
		return result, nil
	}

	return result, glg.Errorf("'%s' uses an unsupported type, expected string or list of strings.", name)
}

// runCommandList takes a list of commands and runs them.
// If a command is prefixed with "#", it will run as root.
func runCommandList(runArgs []string) error {
	for _, v := range runArgs {
		rootOpert := strings.HasPrefix(v, "#")
		newArg := strings.TrimPrefix(v, "#")

		params := []string{cfg.Podman.Container.Name, newArg}

		err := pman.Exec(params, cfg.Image.EnvVars, true, rootOpert,
			podman.Attach{
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				PseudoTTY: true,
			}).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// This is used to parse the recursion of commands.
// It also stops the recursion if it detects a loop.
func parseRecursion(v, name, from string) ([]string, error) {
	parsedName := strings.TrimPrefix(v, "!")
	if parsedName == name || strings.Contains(from, parsedName) {
		return []string{}, glg.Errorf("Recursive command '%s' stopped", name)
	}

	// We pass from where we came from and the name of the command we are parsing.
	// It's also useful for debugging loops when one happens.
	return getAllAsArray(parsedName, fmt.Sprintf("%s,%s", from, name))
}
