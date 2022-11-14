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

// Package state manages the state of the container
package state

import (
	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kadmuffin/develbox/pkg/podman"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

var (
	// Stop is the cobra command for the stop command
	Stop = &cobra.Command{
		Use:     "stop",
		Aliases: []string{"down"},
		Short:   "Stops the container",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Read()
			if err != nil {
				glg.Fatal(err)
			}
			pman := podman.New(cfg.Podman.Path)

			if !pman.Exists(cfg.Podman.Container.Name) {
				glg.Fatal("Container does not exist")
			}

			err = pman.Stop([]string{cfg.Podman.Container.Name}, podman.Attach{})
			if err != nil {
				glg.Fatal(err)
			}
		},
	}
)
