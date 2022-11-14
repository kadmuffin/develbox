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

import "strings"

// ParseArguments splits packages from flags in string arrays. Useful for parsing the packages & flags from Cobra.
func ParseArguments(arguments []string) (packages, flags []string) {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			continue
		}
		packages = append(packages, arg)
	}

	return packages, flags
}

// ContainsString loops through the list and returns true if an item matches the string. Returns true if the item matched our string. False if not.
func ContainsString(list []string, match string) bool {
	for _, item := range list {
		if item == match {
			return true
		}
	}
	return false
}

// RemoveDuplicates returns a new array based on `baseList` without anything repeated on the append list.
func RemoveDuplicates(appendList *[]string, baseList *[]string) []string {
	newList := []string{}
	for _, item := range *baseList {
		if !ContainsString(*appendList, item) {
			newList = append(newList, item)
		}
	}
	return newList
}

// processPackages combines the packages with the argModifier (for example, by adding the prefix or suffix) and returns a string
func processPackages(packages []string, argModifier string) string {
	pkgs := []string{}

	if !strings.Contains(argModifier, "{package}") {
		return strings.Join(packages, " ")
	}

	for _, pkg := range packages {
		pkgs = append(pkgs, strings.ReplaceAll(argModifier, "{package}", pkg))
	}
	return strings.Join(pkgs, " ")
}
