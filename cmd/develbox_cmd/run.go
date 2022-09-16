package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run ...",
	Short: "Run the specified command in the container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}

		develbox.StartContainer(configs.Podman)
		develbox.RunCommands(configs.Commands[args[0]], configs.Podman, true, false)
	},
}

func init() {
	rootCli.AddCommand(run)
}
