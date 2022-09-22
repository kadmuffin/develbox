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

// The pipes package creates a FIFO pipe and provides functions to read/write to it.
//
// After creating the Pipe instance with pipes.New() use pipe.Create() to actually create the pipe file.
//
// It probably only works on Linux, can't confirm if it will work on macOS.
package pipes

import (
	"os"
	"syscall"
)

type Pipe struct {
	path string
}

// New returns a pipe instance.
//
// The path argument is the place where the FIFO pipe will be created.
//
// This constructor is necessary as the path has to be previously defined before  creating a new pipe to mantain the pipe file path consistent.
func New(path string) *Pipe {
	return &Pipe{path: path}
}

// Creates a new FIFO pipe using syscall.Mkfifo(). In case of error the function returns an error
func (e *Pipe) Create() error {
	return syscall.Mkfifo(e.path, 0666)
}

// Removes the pipe file. In case of failure it returns
// the error provided by os.Remove()
func (e *Pipe) Remove() error {
	return os.Remove(e.path)
}

// Waits and reads the contents of the pipe.
/* Use after pipe.Create() */
func (e *Pipe) ReadPipe() ([]byte, error) {
	return os.ReadFile(e.path)
}

// Writes data to the pipe. In case of failure it returns an error.
/* Use after pipe.Create() */
func (e *Pipe) WritePipe(data []byte) error {
	return os.WriteFile(e.path, data, 0666)
}
