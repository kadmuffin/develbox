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

// Package version has the command to print the version of the program
package version

import (
	"github.com/spf13/cobra"
)

var (
	// Number is the current version of the program
	Number = "0.5.4"

	// VersionCmd is the command for printing the current version
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the current version of develbox",
		Long:  `Prints the current version of develbox`,
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			root.Print("v" + Number)
		},
	}
)
