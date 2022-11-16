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

package podman

import (
	"os"
	"os/exec"
	"strings"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
)

// InsideContainer checks if we are inside a container To do this we check the /run/ directory for .containerenv (podman) or .dockerenv (docker)
func InsideContainer() bool {
	inCtnr := config.FileExists("/run/.containerenv") || config.FileExists("/.dockerenv")

	if os.Getenv("CODESPACES") == "true" {
		inCtnr = false
	}
	glg.Debugf("Inside container: %v", inCtnr)
	return inCtnr
}

// ReplaceEnvVars replaces some key environment variables with their values
func ReplaceEnvVars(s string) string {
	// $USER -> os.Getenv("USER")
	//
	// $HOME -> /home/$USER
	//
	// $PWD -> os.Getwd()
	s = strings.ReplaceAll(s, "$$USER", os.Getenv("USER"))
	s = strings.ReplaceAll(s, "$$HOME", "/home/"+os.Getenv("USER"))
	s = strings.ReplaceAll(s, "$$PWD", os.Getenv("PWD"))
	return s
}

// PrintCommand prints the command to run in a more readable format
func PrintCommand(msg string, cmd *exec.Cmd) {
	// - Each argument is on a new line
	//
	// - Flags are on the same line as their argument
	//
	// - The args[1] is on the same line as the command
	//
	// - The command is prefixed with a message
	var args []string
	for i, arg := range cmd.Args {

		if i == 0 {
			args = append(args, arg)
			continue
		}

		if strings.HasPrefix(arg, "-") {
			args = append(args, arg)
		} else {
			args[len(args)-1] = args[len(args)-1] + " " + arg
		}
	}

	glg.Infof(msg, strings.Join(args, "\n  > "))

	// Also print full command
	glg.Infof("Full command: %s", cmd.String())
}

// PrintCommandR prints the command to run in a more readable format and returns the command to run. Format is the same as PrintCommand()
func PrintCommandR(msg string, cmd *exec.Cmd) *exec.Cmd {
	PrintCommand(msg, cmd)
	return cmd
}
