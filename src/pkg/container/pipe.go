package container

import (
	"github.com/kadmuffin/develbox/src/pkg/config"
	"github.com/kadmuffin/develbox/src/pkg/pipes"
	"github.com/kadmuffin/develbox/src/pkg/pkgm"
	"github.com/kpango/glg"
)

// Reads a pipe and installs the requested
// packages.
func readPipe(cfg *config.Struct, pipe *pipes.Pipe) {
	for config.FileExists(pipe.Path()) {
       data, err := pipe.ReadPipe()
	   if err != nil {
		pipe.Remove()
		glg.Errorf("can't read pipe, deleting pipe: %w", err)
	   }

	   opertn, err := pkgm.Read(data)
	   if err != nil {
		glg.Errorf("can't parse json data: %w", err)
	   }

	   opertn.Process(cfg)
	}
}