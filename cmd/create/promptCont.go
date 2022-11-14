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

package create

import (
	"github.com/manifoldco/promptui"
)

// promptEditConfig prompts for no/yes if the user wants to create the container now.
func promptCont() bool {
	prompt := promptui.Prompt{
		Label:     "Create container now",
		IsConfirm: true,
	}

	result, _ := prompt.Run()

	return result == "y"
}
