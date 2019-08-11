package spec

import (
	"strings"

	"golang.org/x/sys/unix"
)

const Version = "1.0.0"

type Spec struct {
	Root     Root    `json:"root"`
	Mounts   []Mount `json:"mounts"`
	Process  Process `json:"process"`
	Hostname string  `json:"hostname"`
	Linux    Linux   `json:"linux"`
}

type Root struct {
	Path     string `json:"path"`
	ReadOnly bool   `json:"readonly"`
}

type Mount struct {
	Destination string   `json:"destination"`
	MountType   string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options"`
}

type Process struct {
	CWD  string   `json:"cwd"`
	Env  []string `json:"env"`
	Args []string `json:"args"`
}

func (mount *Mount) GetFlags(mountType string) uintptr {
	flags := uintptr(0)
	if mountType == "bind" {
		flags |= unix.MS_BIND
	}
	return flags
}

func (mount *Mount) GetOptions(options []string) string {
	return strings.Join(options, ",")
}
