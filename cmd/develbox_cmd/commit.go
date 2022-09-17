package develbox_cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	commit = &cobra.Command{
		Use:   "commit ...",
		Short: "Commits an image.",
		Long: `Commits an image using podman commit.

Any extra arguments will get passed down to podman.
		`,
		PreRun: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StartContainer(configs.Podman)
			command := exec.Command(configs.Podman.Path, append([]string{"commit", configs.Podman.Container.Name}, args...)...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Run()
			develbox.StopContainer(configs.Podman)
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCli.AddCommand(commit)
}
