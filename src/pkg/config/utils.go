package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/kpango/glg"
)

// Checks whether the config folder exists.
//
// Wrapper around config.FileExists()
func ConfigFolderExists() bool {
	return FileExists(".develbox")
}

// Checks whether the config file exists.
//
// Wrapper around config.FileExists()
func ConfigExists() bool {
	return FileExists(".develbox/config.json")
}

// Gets the current folder's full path.
//
// Throws a fatal error in case of failure
// and exists the program.
func GetCurrentDirectory() string {
	currentDir, err := os.Getwd()

	if err != nil {
		glg.Fatalf("failed to get current directory:\n	%w", err)
	}

	return currentDir
}

// Checks if a file/path exists.
// Returns true if it exists
//
// Wrapper around os.Stat() & os.IsNotExists()
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns a hash string made using the current directory's name.
func GetDirNmHash() string {
	currentDirName := filepath.Base(GetCurrentDirectory())
	hasher := sha256.New()
	hasher.Write([]byte(currentDirName))
	dir := hasher.Sum(nil)
	return hex.EncodeToString(dir)

}
