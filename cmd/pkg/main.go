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

// Package pkg has the logic for communication with the package manager
package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/pkgm"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kadmuffin/develbox/pkg/socket"
	"github.com/kpango/glg"
)

// Flags is a struct to hold CLI flags
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

// SendOperation sends an operation to the socket server
func SendOperation(opertn pkgm.Operation) {

	fmt.Println("[Experimental Feature] You *will* need to press enter to continue when the operation is done.")

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
	go func() {
		io.Copy(os.Stdout, reader)
		io.Copy(os.Stderr, errReader)
	}()

	io.Copy(writer, ReadStdinAndWarn(5))
}

// StartContainer starts the container, if we are not inside it.
func StartContainer(cfg *config.Struct) {
	if !podman.InsideContainer() {
		pman := podman.New(cfg.Podman.Path)
		if !pman.Exists(cfg.Container.Name) {
			glg.Fatal("Container does not exist")
		}

		pman.Start([]string{cfg.Container.Name}, podman.Attach{})
	}
}

// ReadStdin reads from stdin and returns the result
func ReadStdin(timeout int) io.Reader {
	// Create a pipe to read from
	reader, writer := io.Pipe()

	// Create a timer to timeout the read
	timer := time.NewTimer(time.Duration(timeout) * time.Second)

	// Create a goroutine to read from stdin
	go func() {
		// Read from stdin
		_, err := io.Copy(writer, os.Stdin)
		if err != nil {
			glg.Error(err)
		}
		// Close the writer
		writer.Close()
	}()

	// Create a goroutine to wait for the timer
	go func() {
		// Wait for the timer to expire
		<-timer.C
		// Close the writer
		writer.Close()
	}()

	// Return the reader
	return reader
}

// ReadStdinAndWarn reads from stdin and warns the user if it takes too long
func ReadStdinAndWarn(timeout int) io.Reader {
	// Create a pipe to read from
	reader, writer := io.Pipe()

	// Create a timer to timeout the read
	timerH := time.NewTimer(time.Duration(timeout) * time.Second)

	// Create a goroutine to read from stdin
	go func() {
		// Read from stdin
		_, err := io.Copy(writer, os.Stdin)
		if err != nil {
			glg.Error(err)
		}
		// Close the writer
		writer.Close()
	}()

	// Create a goroutine to wait for the timer
	go func() {
		<-timerH.C

		fmt.Println("You may want to press enter to continue.")

		writer.Close()
	}()

	// Return the reader
	return reader
}
