package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
)

func main() {
	commands := []string{"touch /test/proven.md"}
	sdir, e := os.Getwd()
	fmt.Println(sdir, e)
	fmt.Println(commands)
	dir := html.EscapeString(filepath.Base(sdir))
	fmt.Println(dir, e)
	createContainer(Image{URI: "debian:testing"}, commands, Podman{Path: "podman", Container: Container{Name: dir, MountBind: "/test:Z"}})
}
