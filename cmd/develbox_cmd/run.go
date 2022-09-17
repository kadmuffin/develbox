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
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands(configs.Commands[args[0]], configs.Podman, true, false, true, false)
		},
	}

	runc = &cobra.Command{
		Use:                "runc ...",
		Short:              "Run the custom command passed in the argument",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommand(args, configs.Podman, true, false, false, "%s", false)
			os.Exit(0)
		},
	}
)

func init() {
	rootCli.AddCommand(run)
	rootCli.AddCommand(runc)
}
