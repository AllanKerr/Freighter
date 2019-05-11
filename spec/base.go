package spec

type Spec struct {
	Root     Root
	Mounts   []Mount
	Process  Process
	Hostname string
	Linux    Linux
}

type Root struct {
	Path     string
	ReadOnly bool
}

type Mount struct {
	Destination string
	MountType   string
	Source      string
	Options     []string
}

type Process struct {
	CWD  string
	Env  []string
	Args []string
}
