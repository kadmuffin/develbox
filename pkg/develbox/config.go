package develbox

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func parseJson(bytes []byte) DevSetings {
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
		hasher.Write([]byte(filepath.Base(getCurrentDirectory())))
		dir := hasher.Sum(nil)
		config.Podman.Container.Name = hex.EncodeToString(dir)
	}
}

func ReadConfig() DevSetings {
	data, err := os.ReadFile(".develbox/config.json")
	if err != nil {
		log.Fatalf("Couldn't read the file .develbox/config.json, exited with: %s", err)
	}

	configs := parseJson(data)
	SetContainerName(&configs)

	return configs
}

func WriteConfig(configs *DevSetings) {
	data, _ := json.MarshalIndent(configs, "", "	")

	err := os.WriteFile(".develbox/config.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func ConfigFolderExists() bool {
	_, err := os.Stat(".develbox")

	return !os.IsNotExist(err)
}

func ConfigExists() bool {
	_, err := os.Stat(".develbox/config.json")

	return !os.IsNotExist(err)
}
