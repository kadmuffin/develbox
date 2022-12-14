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

// Package pipes creates a FIFO pipe and provides functions to read/write to it. It probably only works on Linux, can't confirm if it will work on macOS.
package pipes

import (
	"os"
	"syscall"
)

// Pipe is a FIFO pipe.
type Pipe struct {
	path string
}

// New returns a pipe instance. The path argument is the place where the FIFO pipe will be created.
func New(path string) *Pipe {
	// This constructor is necessary as the path has to be previously defined before  creating a new pipe to mantain the pipe file path consistent.
	return &Pipe{path: path}
}

// Path returns the path to the pipe file.
func (e *Pipe) Path() string {
	return e.path
}

// Create creates a new FIFO pipe using syscall.Mkfifo(). In case of error the function returns an error
func (e *Pipe) Create() error {
	return syscall.Mkfifo(e.path, 0666)
}

// Remove removes the pipe file. In case of failure it returns the error provided by os.Remove()
func (e *Pipe) Remove() error {
	return os.Remove(e.path)
}

// Read waits and reads the contents of the pipe. Use after pipe.Create()
func (e *Pipe) Read() ([]byte, error) {
	return os.ReadFile(e.path)
}

// Write writes data to the pipe. In case of failure it returns an error. Use after pipe.Create()
func (e *Pipe) Write(data []byte) error {
	return os.WriteFile(e.path, data, 0666)
}

// Writer returns a writer to the pipe file.
func (e *Pipe) Writer() (*os.File, error) {
	return os.OpenFile(e.path, os.O_WRONLY, 0666)
}

// Reader returns a reader to the pipe file.
func (e *Pipe) Reader() (*os.File, error) {
	return os.OpenFile(e.path, os.O_RDONLY, 0666)
}

// Exists is wrapper around os.stat() to check if the pipe file exists.
func (e *Pipe) Exists() bool {
	_, err := os.Stat(e.path)
	return !os.IsNotExist(err)
}
