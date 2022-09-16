package develbox

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Create a container with the provided image and name using podman and run the commands
func CreateContainer(config *DevSetings, forceCreation bool) {
	if ContainerExists(config) {
		if !forceCreation {
			log.Fatal("Container already exist!\n ∟ You can force a reset with the -f flag.")
		}
		RemoveContainer(config.Podman)
	}

	fmt.Printf("> Creating container using %s\n", config.Podman.Path)

	currentDir := GetCurrentDirectory()

	cmd := exec.Command(config.Podman.Path, "run", "-a", "stdout", "-a", "stderr", "-it", "-d", config.Podman.Container.Args, "--name", config.Podman.Container.Name, "-v", currentDir+":"+config.Podman.Container.MountPoint, config.Image.URI)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		existCont := ContainerExists(config)
		RemoveContainer(config.Podman)
		if !existCont {
			log.Fatalf("Failed to create container.\n ∟ Command: %s\n%s", cmd.String(), err)
		}
		log.Fatalf("Container started with errors.\n ∟ Command: %s\n%s", cmd.String(), err)
	}

	setupCommands(config)

	if config.Podman.BuildOnly {
		RemoveContainer(config.Podman)
	}
	StopContainer(config.Podman)

	fmt.Print("\nOperation completed.")
}

func StartContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "container", "start", podman.Container.Name)
	cmd.Run()
}

// Remove a container with the provided name using podman
func RemoveContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "rm", "-f", podman.Container.Name)
	cmd.Run()
}

// Stops a container with the provided name using podman
func StopContainer(podman Podman) {
	fmt.Println("\nStopping container...")
	cmd := exec.Command(podman.Path, "stop", podman.Container.Name)
	cmd.Run()
}

// Get the IP address of a container with the provided name using docker
func GetContainerIP(config DevSetings) string {
	cmd := exec.Command(config.Podman.Path, "inspect", config.Podman.Container.Name)
	out, _ := cmd.Output()
	var containerInfo []map[string]interface{}
	json.Unmarshal(out, &containerInfo)
	return containerInfo[0]["NetworkSettings"].(map[string]interface{})["IPAddress"].(string)
}

func setupCommands(config *DevSetings) {
	commands := config.Image.OnCreation

	for _, pkg := range config.Packages {
		commands = append(commands, strings.Replace(config.Image.Installer.Add, "{args}", pkg, 1))
	}

	commands = append(commands, config.Image.OnFinish...)

	RunCommands(commands, config.Podman, true, true)
}

func RunCommands(commandList []string, podman Podman, printOut bool, deleteContainer bool) {
	for _, command := range commandList {
		if len(command) > 0 {
			cmd := exec.Command(podman.Path, "exec", "-it", podman.Container.Name, "sh", "-c", command)
			if printOut {
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			fmt.Printf("Running command: %s\n", cmd.String())
			err := cmd.Run()
			StopContainer(podman)
			if err != nil {
				if deleteContainer {
					RemoveContainer(podman)
				}
				log.Fatal(err)
			}
		}
	}
}

func EnterContainer(config *DevSetings) {
	shell := exec.Command(config.Podman.Path, "exec", "-it", config.Podman.Container.Name, config.Podman.Container.Shell)
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Run()
	StopContainer(config.Podman)
}

func ContainerExists(config *DevSetings) bool {
	cmd := exec.Command(config.Podman.Path, "container", "exists", config.Podman.Container.Name)
	if cmd.Run() == nil {
		return true
	}
	return false
}
