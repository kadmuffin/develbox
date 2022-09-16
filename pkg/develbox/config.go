package develbox

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func parse(bytes []byte) DevSetings {
	var configs DevSetings
	err := json.Unmarshal(bytes, &configs)

	if err != nil {
		log.Fatalf("Couldn't parse the config file, exited with: %s", err)
	}

	SetContainerName(&configs)

	return configs
}

func SetContainerName(config *DevSetings) {
	if config.Podman.Container.Name == "" {
		hasher := sha256.New()
		hasher.Write([]byte(filepath.Base(GetCurrentDirectory())))
		dir := hasher.Sum(nil)
		config.Podman.Container.Name = hex.EncodeToString(dir)
	}
}

func ReadConfig(filename string) DevSetings {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Couldn't read the file %s, exited with: %s", filename, err)
	}

	configs := parse(data)
	SetContainerName(&configs)

	return configs
}

func WriteConfig(configs *DevSetings) {
	data, _ := json.MarshalIndent(configs, "", "	")

	err := os.WriteFile("develbox.json", data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func ConfigExists() bool {
	_, err := os.Stat("develbox.json")

	return !os.IsNotExist(err)
}
