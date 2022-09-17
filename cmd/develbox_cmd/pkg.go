package develbox_cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/kadmuffin/develbox/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	forcePkg bool
	passYes  bool
	pkg      = &cobra.Command{
		Use:   "pkg [add | del | dup | upd]",
		Short: "Manages packages in the container and keeps tracks of them",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	pkg_add = &cobra.Command{
		Use:   "add ...",
		Short: "Installs packages using the pkg manager defined in the config file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			flags := []string{}
			packages := []string{}
			for _, pkg := range args {
				if containsString(configs.Packages, pkg) && !forcePkg {
					fmt.Printf("Skipping installed package: %s (according to develbox.json)", pkg)
					continue
				}

				if string(pkg[0]) == "-" {
					flags = append(flags, pkg)
					continue
				}

				packages = append(packages, pkg)
			}

			pkgM := strings.Replace(configs.Image.Installer.Add, " {-y}", "", 1)
			if passYes {
				pkgM = strings.Replace(configs.Image.Installer.Add, " {-y}", " -y", 1)
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{strings.Replace(pkgM, "{args}", strings.Join(append(flags, packages...), " "), 1)}, configs.Podman, true, false, true, true)
			configs.Packages = append(configs.Packages, packages...)
			develbox.WriteConfig(&configs)
		},
	}

	pkg_del = &cobra.Command{
		Use:   "del ...",
		Short: "Removes packages using the pkg manager defined in the config file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			flags := []string{}
			packages := []string{}
			for _, pkg := range args {
				if !containsString(configs.Packages, pkg) && !forcePkg {
					fmt.Printf("Skipping not installed package: %s (according to develbox.json)", pkg)
					continue
				}

				if string(pkg[0]) == "-" {
					flags = append(flags, pkg)
					continue
				}

				packages = append(packages, pkg)
			}

			pkgM := strings.Replace(configs.Image.Installer.Del, " {-y}", "", 1)
			if passYes {
				pkgM = strings.Replace(configs.Image.Installer.Del, " {-y}", " -y", 1)
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{strings.Replace(pkgM, "{args}", strings.Join(append(flags, packages...), " "), 1)}, configs.Podman, true, false, true, true)
			configs.Packages = append(configs.Packages, packages...)
			develbox.WriteConfig(&configs)
		},
	}

	pkg_upd = &cobra.Command{
		Use:   "upd",
		Short: "Updates the pkg database using the pkg manager defined in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			pkgM := strings.Replace(configs.Image.Installer.Upd, " {-y}", "", 1)
			if passYes {
				pkgM = strings.Replace(configs.Image.Installer.Upd, " {-y}", " -y", 1)
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{pkgM}, configs.Podman, true, false, true, true)
		},
	}

	pkg_dup = &cobra.Command{
		Use:   "dup",
		Short: "Upgrades all packages using the pkg manager defined in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			pkgM := strings.Replace(configs.Image.Installer.Dup, " {-y}", "", 1)
			if passYes {
				pkgM = strings.Replace(configs.Image.Installer.Dup, " {-y}", " -y", 1)
			}

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{pkgM}, configs.Podman, true, false, true, true)
		},
	}
)

func containsString(list []string, match string) bool {
	for _, v := range list {
		if v == match {
			return true
		}
	}
	return false
}

func init() {
	pkg.PersistentFlags().BoolVarP(&forcePkg, "force", "f", false, "Forces the CLI to install/delete packages")
	pkg.PersistentFlags().BoolVarP(&passYes, "yes", "y", false, "Accepts installation before-hand")
	rootCli.AddCommand(pkg)
	pkg.AddCommand(pkg_add)
	pkg.AddCommand(pkg_del)
	pkg.AddCommand(pkg_upd)
	pkg.AddCommand(pkg_dup)

}
