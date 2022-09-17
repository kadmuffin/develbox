package develbox

type Installer struct {
	Add string `default:"apt-get install {args} {-y}" json:"add"` // {-y} needed to auto install on creation
	Del string `default:"apt-get remove {args} {-y}" json:"del"`
	Dup string `default:"apt-get upgrade -y" json:"dup"`
	Upd string `default:"apt-get update" json:"upd"`
}

type Image struct {
	URI        string    `default:"debian:latest" json:"uri"`
	OnCreation []string  `default:"[\"apt-get update\"]" json:"on-creation"`
	OnFinish   []string  `default:"[]" json:"on-finish"`
	Installer  Installer `json:"pkg-manager"`
}

type Binds struct {
	Wayland    bool `json:"wayland"`
	XOrg       bool `json:"xorg"`
	Pulseaudio bool `json:"pulseaudio"`
	Pipewire   bool `json:"pipewire"`
	DRI        bool `json:"dri"`
	Camera     bool `json:"camera"`
}

type Container struct {
	Name     string `json:"name"`
	Args     string `default:"--net=host" json:"arguments"`
	WorkDir  string `default:"/code" json:"work-dir"`
	Shell    string `default:"/bin/bash" json:"shell"`
	RootUser bool   `json:"root-user"`
	Binds    Binds
	Ports    []string `default:"[]" json:"ports"`
	Mounts   []string `default:"[]" json:"mounts"`
}

type Podman struct {
	Path      string    `default:"podman" json:"path"`
	Rootless  bool      `json:"rootless"`
	BuildOnly bool      `json:"build-only"`
	Container Container `json:"container"`
}

type DevSetings struct {
	Image    Image               `json:"image"`
	Podman   Podman              `json:"podman"`
	Commands map[string][]string `default:"{}" json:"commands"`
	Packages []string            `default:"[]" json:"packages"`
}
