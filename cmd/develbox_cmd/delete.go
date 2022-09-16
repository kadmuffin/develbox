package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var delete = &cobra.Command{
	Use:   "delete",
	Short: "Deletes the container if it exists",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}
		develbox.RemoveContainer(configs.Podman)
	},
}

func init() {
	rootCli.AddCommand(delete)
}
