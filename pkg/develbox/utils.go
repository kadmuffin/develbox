package develbox

import (
	"log"
	"os"
)

func GetCurrentDirectory() string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current directory:\n	%s", err)
	}

	return currentDir
}
