package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var root bool

var enter = &cobra.Command{
	Use:   "enter",
	Short: "Enters to the container",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig()
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}
		develbox.StartContainer(configs.Podman)
		develbox.EnterContainer(&configs, root)
	},
}

func init() {
	enter.Flags().BoolVarP(&root, "root", "r", false, "Enters the container with the root user")
	rootCli.AddCommand(enter)
}
