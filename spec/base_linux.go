package spec

import (
	"os"
)

type Linux struct {
	Devices           []Device `json:"devices"`
	RootFSPropagation string   `json:"rootfsPropagation"`
}

type Device struct {
	Path     string      `json:"path"`
	DevType  string      `json:"type"`
	Major    uint32      `json:"major"`
	Minor    uint32      `json:"minor"`
	FileMode os.FileMode `json:"fileMode"`
	UID      int         `json:"uid"`
	GID      int         `json:"gid"`
}
