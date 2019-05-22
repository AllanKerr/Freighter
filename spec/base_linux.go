package spec

import (
	"os"

	"golang.org/x/sys/unix"
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

func (device *Device) GetType() uint32 {
	switch device.DevType {
	case "c", "u":
		return unix.S_IFCHR
	case "b":
		return unix.S_IFBLK
	case "p":
		return unix.S_IFIFO
	}
	return 0
}

func (linux *Linux) GetPropagationType() uintptr {

	switch linux.RootFSPropagation {
	case "slave":
		return unix.MS_SLAVE
	case "rslave":
		return unix.MS_SLAVE | unix.MS_REC
	case "private":
		return unix.MS_PRIVATE
	case "rprivate":
		return unix.MS_PRIVATE | unix.MS_REC
	case "shared":
		return unix.MS_SHARED
	case "rshared":
		return unix.MS_SHARED | unix.MS_REC
	case "unbindable":
		return unix.MS_UNBINDABLE
	case "runbindable":
		return unix.MS_UNBINDABLE | unix.MS_REC
	}
	return 0
}
