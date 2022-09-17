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

	//currentDir := GetCurrentDirectory()

	arguments := []string{"run", "-a", "stdout", "-a", "stderr", "-it", "-d", "--name", config.Podman.Container.Name}

	if config.Podman.Rootless && !config.Podman.Container.RootUser {
		os.Mkdir(".develbox/home", 0755)
		arguments = append(arguments, "--userns=keep-id", "--passwd-entry=develbox:*:$UID:0:develbox_container:/home/develbox:/bin/sh", "-v=./.develbox/home:/home/develbox:Z")
	}

	if len(config.Podman.Container.Args) > 0 {
		arguments = append(arguments, config.Podman.Container.Args)
	}

	if len(config.Podman.Container.Ports) > 0 {
		arguments = append(arguments, processPorts(config.Podman.Container))
	}

	if len(config.Podman.Container.Mounts) > 0 {
		arguments = append(arguments, processVolumes(config.Podman.Container))
	}

	workdir := config.Podman.Container.WorkDir
	if config.Podman.Rootless {
		workdir += ":Z"
	}

	arguments = append(arguments, "-v=.:"+workdir, config.Image.URI)

	cmd := exec.Command(config.Podman.Path, arguments...)
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

	fmt.Print("\nOperation completed.")
}

func StartContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "container", "start", podman.Container.Name)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Container is still shutting down or doesn't exist. %s", err)
	}
}

// Remove a container with the provided name using podman
func RemoveContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "rm", "-f", podman.Container.Name)
	cmd.Run()
}

// Stops a container with the provided name using podman
func StopContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "stop", podman.Container.Name, "-t", "4")
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

	RunCommands(config.Image.OnCreation, config.Podman, true, true, true, true)

	commands := []string{}
	if len(config.Packages) > 0 {
		commands = []string{strings.Replace(strings.Replace(config.Image.Installer.Add, " {-y}", " -y", 1), "{args}", strings.Join(config.Packages, " "), 1)}
	}

	commands = append(commands, config.Image.OnFinish...)

	RunCommands(commands, config.Podman, true, true, false, true)
}

func RunCommands(commandList []string, podman Podman, printOut bool, deleteOnFailure bool, externalStop bool, rootOperation bool) {
	for _, command := range commandList {
		RunCommand([]string{command}, podman, printOut, deleteOnFailure, true, "%s", rootOperation)
	}
	if externalStop {
		return
	}
	StopContainer(podman)
}

func RunCommand(command []string, podman Podman, printOut bool, deleteContainer bool, externalStop bool, errorMessage string, rootOperation bool) []byte {
	shellArgs := []string{"exec", "-it"}
	if podman.Rootless && !podman.Container.RootUser && (string(command[0][0]) == "!" || rootOperation) {
		shellArgs = append(shellArgs, "--user=0")
	}
	shellArgs = append(shellArgs, "-w", podman.Container.WorkDir, podman.Container.Name, "sh", "-c")
	shellArgs = append(shellArgs, command...)

	if string(command[0][0]) != "!" {
		command[0] = strings.Replace(command[0], "!", "", 1)
	}
	if len(command) > 0 {
		cmd := exec.Command(podman.Path, shellArgs...)
		cmd.Stderr = os.Stderr
		var bytes []byte
		var err error
		if printOut {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			fmt.Printf("Running command: %s\n", cmd.String())
			err = cmd.Run()
		} else {
			bytes, err = cmd.Output()
		}
		if !externalStop {
			StopContainer(podman)
		}
		if err != nil {
			if deleteContainer {
				RemoveContainer(podman)
			}
			log.Fatalf(errorMessage, err)
		}
		return bytes
	}
	return nil
}

/*
func Unshare(config DevSetings) {
	if config.Podman.Rootless && (config.Podman.Container.User != "root" && config.Podman.Container.User != "") {
		uid := RunCommand([]string{"id -u " + config.Podman.Container.User}, config.Podman, false, true, true, "Couldn't get the UID of the user "+config.Podman.Container.User+", exited with: %s")

		gid := RunCommand([]string{"id -g " + config.Podman.Container.User}, config.Podman, false, true, true, "Couldn't get the GID of the user "+config.Podman.Container.User+", exited with: %s")

		cmd := exec.Command(config.Podman.Path, "unshare", "chown", strings.ReplaceAll(string(uid)+":"+string(gid), "\r\n", ""), "-R", GetCurrentDirectory())
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			RemoveContainer(config.Podman)
			log.Fatalf("Failed to unshare current folder. %s", err)
		}
	}
}*/

func EnterContainer(config *DevSetings, rootUser bool) {
	shellArgs := []string{"exec", "-it", "-w", config.Podman.Container.WorkDir}
	if rootUser {
		shellArgs = append(shellArgs, "--user=0")
	}
	shellArgs = append(shellArgs, config.Podman.Container.Name, config.Podman.Container.Shell)
	shell := exec.Command(config.Podman.Path, shellArgs...)
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Run()
	/*cmd := exec.Command(config.Podman.Path, "stop", config.Podman.Container.Name, "-t", "2s")
	cmd.Start()*/
}

func ContainerExists(config *DevSetings) bool {
	cmd := exec.Command(config.Podman.Path, "container", "exists", config.Podman.Container.Name)
	return cmd.Run() == nil
}

func processPorts(container Container) string {
	if len(container.Ports) == 0 {
		return ""
	}
	return "-p=" + strings.Join(container.Ports, "-p=")
}

func processVolumes(container Container) string {
	if len(container.Mounts) == 0 {
		return ""
	}
	return "-v=" + strings.ReplaceAll(strings.Join(container.Mounts, "-v="), "{CURRENT_DIR}", getCurrentDirectory())
}
