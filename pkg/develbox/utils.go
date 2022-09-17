package develbox

import (
	"log"
	"os"
)

func getCurrentDirectory() string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current directory:\n	%s", err)
	}

	return currentDir
}
