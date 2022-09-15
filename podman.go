package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Create a container with the provided image and name using podman and run the commands
func createContainer(image Image, commands []string, podman Podman) {
	currentDir, _ := os.Getwd()
	cmd := exec.Command(podman.Path, "run", "-it", "-d", "--name", podman.Container.Name, "-v", currentDir+":"+podman.Container.MountBind, image.URI)
	containerID, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", containerID)

	for _, command := range commands {
		cmd := exec.Command(podman.Path, "exec", podman.Container.Name, "sh", "-c", command)
		err := cmd.Run()
		if err != nil {
			removeContainer(podman)
			log.Fatal(err)
		}
	}
	removeContainer(podman)
}

// Remove a container with the provided name using docker
func removeContainer(podman Podman) {
	cmd := exec.Command(podman.Path, "rm", "-f", podman.Container.Name)
	cmd.Run()
}

// Get the IP address of a container with the provided name using docker
func getContainerIP(containerName string) string {
	cmd := exec.Command("docker", "inspect", containerName)
	out, _ := cmd.Output()
	var containerInfo []map[string]interface{}
	json.Unmarshal(out, &containerInfo)
	return containerInfo[0]["NetworkSettings"].(map[string]interface{})["IPAddress"].(string)
}

func parse(bytes []byte) DevSetings {
	var configs DevSetings
	json.Unmarshal(bytes, &configs)

	return configs
}
