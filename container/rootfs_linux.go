// +build linux

package container

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/allankerr/freighter/log"
	"github.com/allankerr/freighter/spec"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

type linuxRootFS struct {
	path       string
	isReadOnly bool
	prevWD     *string
}

type mountPoint struct {
	src        string
	dst        string
	isRequired bool
	isReadOnly bool
}

func NewRootFS(rootConfig spec.Root) (RootFS, error) {

	rootPath := path.Clean(rootConfig.Path)
	if path.IsAbs(rootPath) {
		wd, err := unix.Getwd()
		if err != nil {
			return nil, err
		}
		log.Info(wd)
		rootPath = path.Join(wd, rootPath)
	}

	return &linuxRootFS{
		path:       rootPath,
		isReadOnly: rootConfig.ReadOnly,
	}, nil
}

func (root *linuxRootFS) PrepareRoot(rootPropagation string) error {

	fileInfo, err := os.Stat(root.path)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Invalid root path, expected directory: %s", root.path)
	}
	propagationType, err := getPropagationType(rootPropagation)
	if err != nil {
		return err
	}
	if err := unix.Mount("", "/", "", propagationType, ""); err != nil {
		return err
	}
	return unix.Mount(root.path, root.path, "bind", unix.MS_BIND|unix.MS_REC, "")
}

func (root *linuxRootFS) CreateMounts(mounts []spec.Mount) (rerr error) {

	if err := root.setWD(); err != nil {
		return err
	}
	defer func() {
		if err := root.restoreWD(); err != nil {
			rerr = err
		}
	}()
	for _, mount := range mounts {
		if err := root.createMount(mount); err != nil {
			return err
		}
	}
	return nil
}

func (root *linuxRootFS) createMount(mount spec.Mount) error {

	dst := mount.Destination
	if path.IsAbs(dst) {
		dst = path.Join(root.path, path.Clean(dst))
	}
	fileInfo, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return err
		}
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("Invalid mount point, expected directory: %s", mount.Destination)
	}

	options := getMountOptions(mount.Options)
	defaultFlags := getMountFlags(mount.MountType)

	return unix.Mount(mount.Source, dst, mount.MountType, defaultFlags, options)
}

func (root *linuxRootFS) CreateDevices(devices []spec.Device) (rerr error) {
	if err := root.setWD(); err != nil {
		return err
	}
	defer func() {
		if err := root.restoreWD(); err != nil {
			rerr = err
		}
	}()
	oldMask := unix.Umask(0000)

	for _, device := range devices {
		if err := root.createDevice(device); err != nil {
			return err
		}
	}
	unix.Umask(oldMask)

	return nil
}

func (root *linuxRootFS) createDevice(device spec.Device) error {

	deviceMode, err := getDeviceMode(device.DevType)
	if err != nil {
		return err
	}
	path := path.Join(root.path, path.Clean(device.Path))
	mode := deviceMode | uint32(device.FileMode)
	dev := unix.Mkdev(device.Major, device.Minor)

	if err := unix.Mknod(path, mode, int(dev)); err != nil {
		return err
	}
	return unix.Chown(path, device.UID, device.GID)
}

func (root *linuxRootFS) PivotRoot() (rerr error) {

	if err := root.setWD(); err != nil {
		return err
	}
	root.clearWD()

	oldRootTarget, err := createOldRootTarget()
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(oldRootTarget); err != nil {
			rerr = err
		}
	}()
	if err := unix.PivotRoot(".", oldRootTarget); err != nil {
		return err
	}
	return unix.Unmount(oldRootTarget, unix.MNT_DETACH)
}

func (root *linuxRootFS) FinalizeRoot() error {
	if root.isReadOnly {
		return unix.Mount("", "/", "", unix.MS_REMOUNT|unix.MS_BIND|unix.MS_RDONLY, "")
	}
	return nil
}

func createOldRootTarget() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	dirName := id.String()
	if err := os.Mkdir(dirName, os.ModePerm); err != nil {
		return "", err
	}
	return dirName, nil
}

func (root *linuxRootFS) setWD() error {

	if root.prevWD != nil {
		return fmt.Errorf("setWD has already been called")
	}
	prevWD, err := os.Getwd()
	if err != nil {
		return err
	}
	root.prevWD = &prevWD
	return os.Chdir(root.path)
}

func (root *linuxRootFS) clearWD() {
	root.prevWD = nil
}

func (root *linuxRootFS) restoreWD() error {

	if err := os.Chdir(*root.prevWD); err != nil {
		return err
	}
	root.prevWD = nil
	return nil
}

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
