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
	"regexp"
	"strings"

	"github.com/kadmuffin/develbox/cmd/state"
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	cfg  config.Structure
	pman podman.Podman

	// Run is the command to run command defined in config.
	Run = &cobra.Command{
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
			if !pman.Exists(cfg.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			state.StartContainer(cfg.Container.Name, pman, podman.Attach{})

			name := strings.Join(args, " ")
			if _, ok := cfg.Commands[name]; !ok {
				return glg.Errorf("Command '%s' does not exist", name)
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
func getAllAsArray(name string, from string) ([]string, error) {
	// The resulting array will contain all the commands inside that name.
	//
	// If the name is prefixed with "!", it will recursively call itself
	// to get the full command tree.
	if _, ok := cfg.Commands[name]; !ok {
		return []string{}, glg.Errorf("[%s] Command '%s' does not exist", from, name)
	}
	cmds := cfg.Commands[name]
	result := []string{}

	switch cmds := cmds.(type) {
	case string:
		if strings.HasPrefix(cmds, "!") {
			parsedCmds, err := runBashParse(cmds)
			if err != nil {
				return []string{}, err
			}
			return parseRecursion(parsedCmds, name, from)
		}
		result = append(result, cmds)
		return result, nil
	case []interface{}:
		for _, v := range cmds {
			if strings.HasPrefix(v.(string), "!") {
				parsedCmds, err := runBashParse(v.(string))
				if err != nil {
					return []string{}, err
				}

				newCmds, err := parseRecursion(parsedCmds, name, from)

				if err != nil {
					return []string{}, err
				}

				result = append(result, newCmds...)
				continue
			}
			result = append(result, v.(string))
		}
		return result, nil
	default:
		return result, glg.Errorf("'%s' uses an unsupported type, expected string or list of strings.", name)
	}

}

// runCommandList takes a list of commands and runs them. If a command is prefixed with "#", it will run as root.
func runCommandList(runArgs []string) error {
	for _, v := range runArgs {
		rootOpert := strings.HasPrefix(v, "#")
		newArg := strings.TrimPrefix(v, "#")

		params := []string{cfg.Container.Name, newArg}

		err := pman.Exec(params, cfg.Image.Variables, true, rootOpert,
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

// parseRecursion is used to parse the recursion of commands. It also stops the recursion if it detects a loop.
func parseRecursion(v, name, from string) ([]string, error) {
	parsedName := strings.TrimPrefix(v, "!")
	if parsedName == name || strings.Contains(from, parsedName) {
		return []string{}, glg.Errorf("Recursive command '%s' stopped", name)
	}

	// We pass from where we came from and the name of the command we are parsing.
	// It's also useful for debugging loops when one happens.
	return getAllAsArray(parsedName, fmt.Sprintf("%s,%s", from, name))
}

// parseSubBash runs commands inside "${}" and returns the result.
//
// (It runs them inside the container, then it auto replaces itself with the result)
//
// Matches in any place of the string. (Runs all matches)
func parseSubBash(v string) (string, error) {
	re := regexp.MustCompile(`\$\{(.+?)\}`)
	matches := re.FindAllStringSubmatch(v, -1)

	for _, match := range matches {
		params := []string{cfg.Container.Name, match[1]}
		result, err := pman.Exec(params, cfg.Image.Variables, true, false, podman.Attach{}).Output()
		if err != nil {
			return "", err
		}

		v = strings.Replace(v, match[0], strings.TrimSuffix(string(result), "\n"), -1)
	}

	return v, nil
}

// parseSubRootBash runs commands inside "$#{}" and returns the result.
//
// (It runs them inside the container, then it auto replaces itself with the result)
//
// Matches in any place of the string. (Runs all matches)
func parseSubRootBash(v string) (string, error) {
	re := regexp.MustCompile(`\$#\{(.+?)\}`)
	matches := re.FindAllStringSubmatch(v, -1)

	for _, match := range matches {
		params := []string{cfg.Container.Name, match[1]}
		result, err := pman.Exec(params, cfg.Image.Variables, true, true, podman.Attach{}).Output()
		if err != nil {
			return "", err
		}

		v = strings.Replace(v, match[0], strings.TrimSuffix(string(result), "\n"), -1)
	}

	return v, nil
}

// runBashParse runs the bash parse and returns the result.
func runBashParse(v string) (string, error) {
	v, err := parseSubBash(v)
	if err != nil {
		return "", err
	}

	v, err = parseSubRootBash(v)
	if err != nil {
		return "", err
	}

	return v, nil
}
