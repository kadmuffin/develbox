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

package pkgm

func ParseArguments(arguments []string) (*[]string, *[]string) {
	packages := []string{}
	flags := []string{}
	for _, arg := range arguments {
		if string(arg[0]) == "-" {
			flags = append(flags, arg)
			continue
		}
		packages = append(packages, arguments...)
	}

	return &packages, &flags
}

func ContainsString(list []string, match string) bool {
	for _, item := range list {
		if item == match {
			return true
		}
	}
	return false
}

func RemoveDuplicates(baseList *[]string, appendList *[]string) []string {
	newList := []string{}
	for _, item := range *appendList {
		if !ContainsString(*baseList, item) {
			newList = append(newList, item)
		}
	}
	return newList
}
