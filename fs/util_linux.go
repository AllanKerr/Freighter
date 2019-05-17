package fs

import (
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

func getDeviceMode(mode string) (uint32, error) {
	switch mode {
	case "c", "u":
		return unix.S_IFCHR, nil
	case "b":
		return unix.S_IFBLK, nil
	case "p":
		return unix.S_IFIFO, nil
	}
	return 0, fmt.Errorf("Unexpected device mode: %s", mode)
}

func getPropagationType(propagationType string) (uintptr, error) {

	if propagationType == "slave" {
		return unix.MS_SLAVE, nil
	}
	if propagationType == "rslave" {
		return unix.MS_SLAVE | unix.MS_REC, nil
	}
	if propagationType == "private" {
		return unix.MS_PRIVATE, nil
	}
	if propagationType == "rprivate" {
		return unix.MS_PRIVATE | unix.MS_REC, nil
	}
	if propagationType == "shared" {
		return unix.MS_SHARED, nil
	}
	if propagationType == "rshared" {
		return unix.MS_SHARED | unix.MS_REC, nil
	}
	if propagationType == "unbindable" {
		return unix.MS_UNBINDABLE, nil
	}
	if propagationType == "runbindable" {
		return unix.MS_UNBINDABLE | unix.MS_REC, nil
	}
	if propagationType == "" {
		return 0, nil
	}
	return 0, fmt.Errorf("Invalid root FS propagation: %s", propagationType)
}

func getMountFlags(mountType string) uintptr {
	if mountType == "bind" {
		return unix.MS_BIND
	}
	return 0
}

func getMountOptions(options []string) string {
	return strings.Join(options, ",")

}
