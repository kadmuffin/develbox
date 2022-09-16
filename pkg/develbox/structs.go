package develbox

type Installer struct {
	Add string `default:"apt-get install {args}" json:"add"`
	Del string `default:"apt-get remove {args}" json:"del"`
	Dup string `default:"apt-get upgrade -y" json:"dup"`
	Upd string `default:"apt-get update" json:"upd"`
}

type Image struct {
	URI        string    `default:"debian:latest" json:"uri"`
	OnCreation []string  `default:"[\"apt-get update\", \"\"]" json:"on-creation"`
	OnFinish   []string  `default:"[]" json:"on-finish"`
	Installer  Installer `json:"pkg-manager"`
}

type Container struct {
	Name       string `json:"name"`
	User       string `default:"root" json:"user"`
	Args       string `default:"--net=host" json:"arguments"`
	MountPoint string `default:"/code:Z" json:"mount-point"`
	Shell      string `default:"/bin/bash" json:"shell"`
}

type Podman struct {
	Path      string    `default:"podman" json:"path"`
	Rootless  bool      `json:"unshare"`
	BuildOnly bool      `json:"build-only"`
	Container Container `json:"container"`
}

type DevSetings struct {
	Image    Image               `json:"image"`
	Podman   Podman              `json:"podman"`
	Commands map[string][]string `default:"{}" json:"commands"`
	Packages []string            `default:"[]" json:"packages"`
}
