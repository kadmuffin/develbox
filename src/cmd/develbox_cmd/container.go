package develbox_cmd

import (
	"log"

	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	container = &cobra.Command{
		Use:   "container",
		Short: "Manage the container state with the sub-commands",
	}

	start = &cobra.Command{
		Use:   "start",
		Short: "Start the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StartContainer(configs.Podman)
		},
	}

	stop = &cobra.Command{
		Use:   "stop",
		Short: "Stop the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StopContainer(configs.Podman)
		},
	}

	delete = &cobra.Command{
		Use:   "delete",
		Short: "Delete the develbox container",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}
			develbox.StopContainer(configs.Podman)
		},
	}
)

func init() {
	container.AddCommand(start)
	container.AddCommand(stop)
	container.AddCommand(delete)

	rootCli.AddCommand(container)
}
