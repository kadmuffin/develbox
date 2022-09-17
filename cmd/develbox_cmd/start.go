package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var start = &cobra.Command{
	Use:   "start",
	Short: "Starts the container if it exists",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig()
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}
		develbox.StartContainer(configs.Podman)
	},
}

func init() {
	rootCli.AddCommand(start)
}
