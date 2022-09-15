package main

type Installer struct {
	Add string `default:"apt-get install %s -y"`
	Del string `default:"apt-get remove %s -y"`
}

type Image struct {
	URI        string `default:"debian:testing"`
	OnCreation []string
	OnFinish   []string
	Installer  Installer
}

type Container struct {
	Name      string
	User      string `default:"root"`
	Args      string
	MountBind string `default:"/dev"`
}

type Podman struct {
	Path      string `default:"podman"`
	Rootless  bool
	Container Container
}

type DevSetings struct {
	Image    Image
	Commands map[string]string
	Podman   Podman
	Packages []string
}
