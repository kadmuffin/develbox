package develbox_cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var forcePkg bool

var pkg = &cobra.Command{
	Use:   "pkg [add | del | dup]",
	Short: "Manages packages in the container and keeps tracks of them",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var pkg_add = &cobra.Command{
	Use:   "add ...",
	Short: "Installs packages using the pkg manager defined in the config file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}

		var packages string = ""
		for _, pkg := range args {
			if containsString(configs.Packages, pkg) && !forcePkg {
				fmt.Printf("Skipping installed package: %s (according to develbox.json)", pkg)
				continue
			}

			packages += pkg + " "
		}

		develbox.StartContainer(configs.Podman)
		develbox.RunCommands([]string{strings.Replace(configs.Image.Installer.Add, "{args}", packages, 1)}, configs.Podman, true, false)
		configs.Packages = append(configs.Packages, args...)
		develbox.WriteConfig(&configs)
	},
}

var pkg_del = &cobra.Command{
	Use:   "del ...",
	Short: "Removes packages using the pkg manager defined in the config file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}

		var packages string = ""
		for _, pkg := range args {
			if !containsString(configs.Packages, pkg) && !forcePkg {
				fmt.Printf("Skipping not installed package: %s (according to develbox.json)", pkg)
				continue
			}
			packages += pkg + " "
		}

		develbox.StartContainer(configs.Podman)
		develbox.RunCommands([]string{strings.Replace(configs.Image.Installer.Del, "{args}", packages, 1)}, configs.Podman, true, false)
		configs.Packages = append(configs.Packages, args...)
		develbox.WriteConfig(&configs)
	},
}

var pkg_upd = &cobra.Command{
	Use:   "upd",
	Short: "Updates the pkg database using the pkg manager defined in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}

		develbox.StartContainer(configs.Podman)
		develbox.RunCommands([]string{configs.Image.Installer.Upd}, configs.Podman, true, false)
	},
}

var pkg_dup = &cobra.Command{
	Use:   "dup",
	Short: "Upgrades all packages using the pkg manager defined in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		var configs develbox.DevSetings = develbox.ReadConfig("develbox.json")
		if !develbox.ContainerExists(&configs) {
			log.Fatal("No container found")
		}

		develbox.StartContainer(configs.Podman)
		develbox.RunCommands([]string{configs.Image.Installer.Dup}, configs.Podman, true, false)
	},
}

func containsString(list []string, match string) bool {
	for _, v := range list {
		if v == match {
			return true
		}
	}
	return false
}

func init() {
	pkg.Flags().BoolVarP(&forcePkg, "force", "f", false, "Forces the CLI to install/delete packages")
	rootCli.AddCommand(pkg)
	pkg.AddCommand(pkg_add)
	pkg.AddCommand(pkg_del)
	pkg.AddCommand(pkg_upd)
	pkg.AddCommand(pkg_dup)

}
