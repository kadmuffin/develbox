package develbox

type Installer struct {
	Add string
	Del string
}

type Image struct {
	Name       string
	OnCreation []string
	OnFinish   []string
	Installer  Installer
}

type Podman struct {
	Path     string
	RunAs    string
	Rootless bool
}

type DevSetings struct {
	Image    Image
	Commands map[string]string
	Podman   Podman
	Packages []string
}
