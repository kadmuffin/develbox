package develbox_cmd

import (
	"log"
	"os"

	"github.com/creasty/defaults"
	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var writeConfig bool

var create = &cobra.Command{
	Use:   "create",
	Short: "Creates a container using the develbox.json",
	Run: func(cmd *cobra.Command, args []string) {

		if writeConfig {
			configs := new(develbox.DevSetings)
			develbox.SetContainerName(configs)
			defaults.Set(configs)
			os.Mkdir(".develbox", 0755)
			if develbox.ConfigExists() && !forceAction {
				log.Fatal("A config file already exist!")
			}
			os.Remove(".develbox/config.json")
			develbox.WriteConfig(configs)
			os.Exit(0)
		}
		var configs develbox.DevSetings = develbox.ReadConfig()
		develbox.CreateContainer(&configs, forceAction)
	},
}

func init() {
	create.Flags().BoolVarP(&forceAction, "force", "f", false, "Forces to overwrite the container/config")
	create.Flags().BoolVarP(&writeConfig, "config", "c", false, "Writes a config file with pre-set defaults")
	rootCli.AddCommand(create)
}
