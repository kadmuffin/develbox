package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var stop = &cobra.Command{
	Use:   "stop",
	Short: "Stops the container if it exists",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}
		develbox.StopContainer(configs.Podman)
	},
}

func init() {
	rootCli.AddCommand(stop)
}
