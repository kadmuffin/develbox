package develbox_cmd

import (
	"github.com/spf13/cobra"
)

var (
	forceAction bool
	rootCli     = &cobra.Command{
		Use:   "develbox",
		Short: "Develbox - Simple CLI tool useful for managing dev enviroments.",
		Long: `Develbox - A simple but dirty CLI tool that manages containerized dev enviroments.

Created so I don't have to expose my entire computer to random node modules (and to learn Go, that means BAD CODE).`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Execute() {
	rootCli.PersistentFlags().BoolVarP(&forceAction, "force", "f", false, "Forces the subsequent action to execute.")
	rootCli.Execute()
}
