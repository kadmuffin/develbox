package develbox_cmd

import (
	"log"
	"os"

	"github.com/creasty/defaults"
	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var forceCreate bool
var writeConfig bool

var create = &cobra.Command{
	Use:   "create",
	Short: "Creates a container using the develbox.json",
	Run: func(cmd *cobra.Command, args []string) {

		if writeConfig {
			configs := new(develbox.DevSetings)
			develbox.SetContainerName(configs)
			defaults.Set(configs)
			if develbox.ConfigExists() && !forceCreate {
				log.Fatal("A config file already exist!")
			}
			os.Remove("develbox.json")
			develbox.WriteConfig(configs)
			os.Exit(0)
		}
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		develbox.CreateContainer(&configs, forceCreate)
	},
}

func init() {
	create.Flags().BoolVarP(&forceCreate, "force", "f", false, "Forces to overwrite the container/config")
	create.Flags().BoolVarP(&writeConfig, "config", "c", false, "Writes a config file with pre-set defaults")
	rootCli.AddCommand(create)
}
