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
