package develbox_cmd

import (
	"log"
	"os"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	customCommand bool

	run = &cobra.Command{
		Use:   "run ...",
		Short: "Run the specified command defined in the config file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			develbox.StartContainer(configs.Podman)
			if customCommand {
				develbox.RunCommand(args, configs.Podman, true, false, false, "%s", false)
				os.Exit(0)
			}
			develbox.RunCommands(configs.Commands[args[0]], configs.Podman, true, false, true, false)
		},
	}
)

func init() {
	run.Flags().BoolVarP(&customCommand, "custom", "c", false, "Runs the arguments directly in the container")
	rootCli.AddCommand(run)
}
