package podman

import "github.com/kadmuffin/develbox/src/pkg/config"

// Checks if we are inside a container
//
// To do this we check the /run directory for .containerenv (podman)
// or .dockerenv (docker)
func InsideContainer() bool {
	return config.FileExists("/run/.containerenv") || config.FileExists("/run/.dockerenv")
}
