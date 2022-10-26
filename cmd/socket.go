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

package cmd

import (
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kadmuffin/develbox/pkg/socket"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	Socket = &cobra.Command{
		Use:   "socket",
		Short: "Creates a socket that enables communication with the container",
		Long: `Creates a socket that enables communication with the container
		
		Used so we can install packages from inside the container (without using root).`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Read()
			if err != nil {
				glg.Failf("Can't read config: %s", err)
			}
			pman := podman.New(cfg.Podman.Path)
			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}
			pman.Start([]string{cfg.Podman.Container.Name}, podman.Attach{})

			defer os.Remove(".develbox/home/.develbox.sock")
			createSocket(&cfg)
		},
	}
)

func createSocket(cfg *config.Struct) {
	// Remove socket file, just in case
	os.Remove(".develbox/home/.develbox.sock")

	// Create a socket to communicate with the container
	s := socket.New(".develbox/home/.develbox.sock")

	if !s.Exists() {
		glg.Debug("Socket doesn't exist, creating...")
		if err := s.Create(); err != nil {
			glg.Fatal(err)
		}
	}

	defer s.Close()

	for {
		glg.Debug("Waiting for requests...")
		s.Listen(func() {
			defer s.CloseConnection()

			// Use socket as tty routing so that
			// inside the container we can pass the tty to the host
			// and run a command with that tty.
			operation := pkgm.Operation{}
			err := s.ReceiveJSON(&operation)
			if err != nil {
				glg.Fatal(err)
			}

			// Print the operation as JSON text
			glg.Debug(operation.ToJSON())
			command, err := operation.ProcessCmd(cfg, podman.Attach{})
			if err != nil {
				glg.Fatal(err)
			}

			reader, _ := s.Reader()
			writer, _ := s.Writer()
			errWriter, _ := s.Writer()

			command.Stdin = reader
			command.Stdout = writer
			command.Stderr = errWriter

			glg.Debug("Running command: ", command)
			err = command.Run()
			if err == nil {
				operation.UpdateConfig(cfg)
				config.Write(cfg)
			}

			glg.Debug("Command finished with error: %v\n", err)
		})
	}
}
