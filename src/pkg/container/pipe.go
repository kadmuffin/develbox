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

package container

import (
	"os"

	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/pipes"
	"github.com/kadmuffin/develbox/src/pkg/pkgm"
	"github.com/kpango/glg"
)

// Reads a pipe and installs the requested
// packages.
func pkgPipe(cfg *config.Struct, pipe *pipes.Pipe) {
	for {
		data, err := pipe.Read()
		if err != nil {
			pipe.Remove()
			glg.Errorf("can't read pipe, deleting pipe: %s", err)
		}

		opertn, err := pkgm.Read(data)
		if err != nil {
			glg.Errorf("can't parse json data: %s", err)
		}

		cmd, err := opertn.ProcessCmd(cfg)
		if err != nil {
			glg.Errorf("operation type is invalid: %s", err)
		}

		file, _ := os.Open(pipe.Path())
		cmd.Stdin = file
		cmd.Run()
		pipe.Remove()
	}
}
