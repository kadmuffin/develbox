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
	"io"
	"os"
	"path/filepath"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kadmuffin/develbox/pkg/socket"
	"github.com/kpango/glg"
)

type Flags struct {
	All       []string
	DevPkg    bool
	UserOpert bool
	ShowHelp  bool
}

func parseFlags(flags *[]string) (result Flags) {
	for _, flag := range *flags {
		switch flag {
		case "--dev":
			result.DevPkg = true
		case "-D":
			result.DevPkg = true
		case "--user":
			result.UserOpert = true
		case "-U":
			result.UserOpert = true
		case "--help":
			result.ShowHelp = true
		case "-h":
			result.ShowHelp = true
		case "--pkg-help":
			result.All = append(result.All, "--help")
		case "-p":
			result.All = append(result.All, "--help")

		default:
			result.All = append(result.All, flag)
		}
	}
	return
}

// Sends an operation to the socket server
func SendOperation(opertn pkgm.Operation) {
	home := os.Getenv("HOME")

	glg.Debugf("Home: %s", home)
	s := socket.New(filepath.Join(home, ".develbox.sock"))

	if !s.Exists() {
		glg.Fatal("Socket does not exist")
	}

	s.Connect()

	// We first pass the operation to the socket
	s.SendJSON(opertn)

	// Then we attach the socket to stdout
	// And attach stdin to the socket
	reader, _ := s.Reader()
	writer, _ := s.Writer()
	errReader, _ := s.Reader()

	// We copy the socket to stdin,stdout,stderr
	// But, we also let the user type into the terminal
	go io.Copy(writer, os.Stdin)
	go io.Copy(os.Stderr, errReader)
	io.Copy(os.Stdout, reader)
}

// Starts the container, if we are not inside it.
func StartContainer(cfg *config.Struct) {
	if !podman.InsideContainer() {
		pman := podman.New(cfg.Podman.Path)
		if !pman.Exists(cfg.Podman.Container.Name) {
			glg.Fatal("Container does not exist")
		}

		pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})
	}
}
