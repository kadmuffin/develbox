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

// Prompts the user about .dockerignore
package dockerfile

import (
	"fmt"
	"os"

	"github.com/kadmuffin/develbox/pkg/config"
	"github.com/kpango/glg"
	"github.com/manifoldco/promptui"
)

// Create a .dockerignore file using the .gitignore file
func createDockerignore() {
	// Create a .dockerignore file
	file, err := os.Create(".dockerignore")
	if err != nil {
		glg.Fatal(err)
	}
	defer file.Close()

	// Copy the contents of .gitignore to .dockerignore
	ignore, err := os.Open(".gitignore")
	if err != nil {
		glg.Fatal(err)
	}

	// Write the contents of .gitignore to .dockerignore
	_, err = file.ReadFrom(ignore)

	if err != nil {
		glg.Fatal(err)
	}

	// Add .git to .dockerignore
	file.WriteString(".git")

	glg.Info("Created .dockerignore file")
}

func selectDck() bool {
	items := []string{}

	if config.FileExists(".dockerignore") {
		items = append(items, "Use .dockerignore (recommended)")
	} else {
		items = append(items, "I will create a .dockerignore (recommended)")
	}

	if config.FileExists(".gitignore") {
		items = append(items, "Create .dockerignore (with .gitignore contents, not ideal)")
		items = append(items, "Layer files using .gitignore (image size will be huge)")
	}

	items = append(items, "Copy all files (not recommended!)")

	prompt := promptui.Select{
		Label: "Choose what to do with .dockerignore",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		glg.Fatal(err)
	}

	switch result {
	case "Use .dockerignore (recommended)":
		return true
	case "I will create a .dockerignore (recommended)":
		fmt.Println("Run this command again after you make your .dockerignore file")
		os.Exit(0)
	case "Create .dockerignore (may not work as expected)":
		createDockerignore()
		return true
	case "Layer files using .gitignore (image size will be huge)":
		return false
	case "Copy all files (not recommended!)":
		return true
	}
	return true
}
