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

package global

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kpango/glg"
)

// Define a regex so only valid folder names are allowed
var validFolderName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString

// Creates a new shared folder at $XDG_DATA_HOME/develbox/shared/<name>
func CreateTaggedFolder(tag string) error {
	if !validFolderName(tag) {
		glg.Fatalf("Invalid tag name, only alphanumeric characters and underscores are allowed. Got %s", tag)
	}

	err := CreateFolder(GetDataHome() + "/develbox/shared/" + tag)
	if err != nil && !os.IsExist(err) {
		glg.Fatal(err)
	}
	return nil
}

// Returns the path to a shared folder at $XDG_DATA_HOME/develbox/shared/<name>
func GetTaggedFolder(tag string) string {
	return GetDataHome() + "/develbox/shared/" + tag
}

// Creates a new shared folder and returns the path to it
func CreateAndGet(tag string) string {
	CreateTaggedFolder(tag)
	return GetTaggedFolder(tag)
}

// Wrapper around HashPath that creates the folder if it doesn't exist
// inside the shared folder
func HashPathAndCreate(path string, tag string) (string, error) {
	hashPath := HashPath(path)
	err := CreateFolder(fmt.Sprintf("%s/%s", GetTaggedFolder(tag), hashPath))
	return hashPath, err
}
