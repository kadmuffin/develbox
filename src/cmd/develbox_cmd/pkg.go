package develbox_cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/kadmuffin/develbox/src/pkg/develbox"
	"github.com/spf13/cobra"
)

var (
	pkg_add = &cobra.Command{
		Use:                "add ...",
		Short:              "Installs packages using the pkg manager defined in the config file",
		DisableFlagParsing: true,
		Args:               cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			packages := args
			packagesForced := []string{}
			flags := []string{}

			for i, pkg := range packages {
				if string(pkg[0]) == "-" {
					if i == 0 && pkg == "-f" {
						forceAction = true
					} else {
						flags = append(flags, pkg)
					}

					packages = append(packages[:i], packages[i+1:]...)
				}
			}

			for i, pkg := range packages {
				if containsString(configs.Packages, pkg) {
					if !forceAction {
						fmt.Printf("Skipping installed package: %s (according to develbox.json)\n", pkg)
						packages = append(packages[:i], packages[i+1:]...)

						continue
					}

					packagesForced = append(packagesForced, pkg)
					continue
				}

			}
			if !forceAction && len(packages) <= 0 {
				fmt.Println("No packages were installed.")
				return
			}
			pkgM := strings.Replace(configs.Image.Installer.Add, " {-y}", "", 1)

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{strings.Replace(pkgM, "{args}", strings.Join(append(flags, append(packages, packagesForced...)...), " "), 1)}, configs.Podman, true, false, false, true, true)
			configs.Packages = append(configs.Packages, packages...)
			develbox.WriteConfig(&configs)
		},
	}

	pkg_del = &cobra.Command{
		Use:                "del ...",
		Short:              "Removes packages using the pkg manager defined in the config file",
		DisableFlagParsing: true,
		Args:               cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			packages := args
			packagesForced := []string{}
			flags := []string{}

			for i, pkg := range packages {
				if string(pkg[0]) == "-" {
					if i == 0 && pkg == "-f" {
						forceAction = true
					} else {
						flags = append(flags, pkg)
					}

					packages = append(packages[:i], packages[i+1:]...)
				}
			}

			for i, pkg := range packages {
				if !containsString(configs.Packages, pkg) {
					if !forceAction {
						fmt.Printf("Skipping not installed package: %s (according to develbox.json)\n", pkg)
						packages = append(packages[:i], packages[i+1:]...)
						continue
					}

					packagesForced = append(packagesForced, pkg)
					continue
				}

			}
			if !forceAction && len(packages) <= 0 {
				fmt.Println("No packages were deleted.")
				return
			}
			pkgM := strings.Replace(configs.Image.Installer.Del, " {-y}", "", 1)

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{strings.Replace(pkgM, "{args}", strings.Join(append(flags, append(packages, packagesForced...)...), " "), 1)}, configs.Podman, true, false, false, true, true)
			for i, pkg := range packages {
				for _, delPkgs := range packages {
					if delPkgs == pkg {
						configs.Packages = append(configs.Packages[:i], configs.Packages[:i+1]...)
					}
				}
			}
			develbox.WriteConfig(&configs)
		},
	}

	pkg_srch = &cobra.Command{
		Use:                "search",
		Short:              "Searches a pkg using the pkg manager defined in the config file",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			pkgM := strings.Replace(configs.Image.Installer.Srch, " {-y}", "", 1)
			pkgM = strings.Replace(pkgM, "{args}", strings.Join(args, " "), 1)

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{pkgM}, configs.Podman, true, false, false, true, true)
		},
	}

	pkg_upd = &cobra.Command{
		Use:                "update",
		Short:              "Updates the pkg database using the pkg manager defined in the config file",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			pkgM := strings.Replace(configs.Image.Installer.Upd, " {-y}", "", 1)

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{pkgM}, configs.Podman, true, false, false, true, true)
		},
	}

	pkg_dup = &cobra.Command{
		Use:                "upgrade",
		Short:              "Upgrades all packages using the pkg manager defined in the config file",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			var configs develbox.DevSetings = develbox.ReadConfig()
			if !develbox.ContainerExists(&configs) {
				log.Fatal("No container found")
			}

			pkgM := strings.Replace(configs.Image.Installer.Dup, " {-y}", "", 1)

			develbox.StartContainer(configs.Podman)
			develbox.RunCommands([]string{pkgM}, configs.Podman, true, false, false, true, true)
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
	rootCli.AddCommand(pkg_add)
	rootCli.AddCommand(pkg_del)
	rootCli.AddCommand(pkg_srch)
	rootCli.AddCommand(pkg_upd)
	rootCli.AddCommand(pkg_dup)
}
