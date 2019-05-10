package fs

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

type linuxRootFS struct {
	target string
}

type mountPoint struct {
	src        string
	dst        string
	isRequired bool
	isReadOnly bool
}

func NewRootFS(target string) (RootFS, error) {

	if _, err := os.Stat(target); err != nil {
		return nil, err
	}
	return &linuxRootFS{target}, nil
}

func (root *linuxRootFS) PrepareRoot() error {

	if err := unix.Mount(root.target, root.target, "bind", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return err
	}
	if err := os.Chdir(root.target); err != nil {
		return err
	}
	if err := mountProc(); err != nil {
		return err
	}
	if err := mountDev(); err != nil {
		return err
	}
	return createDefaultDevices()
}

func mountProc() error {

	if err := os.MkdirAll("proc", 0755); err != nil {
		return err
	}
	return unix.Mount("proc", "proc", "proc", 0, "")
}

func mountDev() error {
	if err := os.Mkdir("dev", 0755); err != nil {
		return err
	}
	return unix.Mount("tmpfs", "dev", "tmpfs", 0, "size=65536k")
}

func createDefaultDevices() error {

	if err := createDevice("dev/null", unix.S_IFCHR, 1, 3, 0666); err != nil {
		return err
	}
	if err := createDevice("dev/zero", unix.S_IFCHR, 1, 5, 0666); err != nil {
		return err
	}
	if err := createDevice("dev/full", unix.S_IFCHR, 1, 7, 0666); err != nil {
		return err
	}
	if err := createDevice("dev/random", unix.S_IFCHR, 1, 8, 0666); err != nil {
		return err
	}
	if err := createDevice("dev/urandom", unix.S_IFCHR, 1, 9, 0666); err != nil {
		return err
	}
	return createDevice("dev/tty", unix.S_IFCHR, 5, 0, 0666)
}

func createDevice(path string, mode uint32, major uint32, minor uint32, fileMode os.FileMode) error {

	dev := unix.Mkdev(major, minor)
	if err := unix.Mknod(path, mode, int(dev)); err != nil {
		return err
	}
	return os.Chmod(path, fileMode)
}

func (root *linuxRootFS) AddSystemCommands() error {

	if err := os.Chdir(root.target); err != nil {
		return err
	}

	// TODO, convert to use overlayfs instead of RO bind mounts
	mountPoints := getSystemCommandMountPoints()
	for _, mountPoint := range mountPoints {
		if err := createReadOnlyMount(mountPoint); err != nil {
			return err
		}
	}
	return nil
}

func getSystemCommandMountPoints() []mountPoint {
	return []mountPoint{
		mountPoint{"/bin", "bin", true, true},
		mountPoint{"/lib", "lib", true, true},
		mountPoint{"/lib32", "lib32", false, true},
		mountPoint{"/lib64", "lib64", false, true},
		mountPoint{"/usr/bin", "usr/bin", true, true},
		mountPoint{"/usr/lib", "usr/lib", true, true},
		mountPoint{"/usr/lib32", "usr/lib32", false, true},
		mountPoint{"/usr/lib64", "usr/lib64", false, true},
	}
}

func createReadOnlyMount(mountPoint mountPoint) error {

	fileInfo, err := os.Stat(mountPoint.src)
	if err != nil {
		if mountPoint.isRequired {
			return err
		}
		return nil
	}
	mode := fileInfo.Mode()
	if (mode&os.ModeDir) == 0 && (mode&os.ModeSymlink) == 0 {
		if mountPoint.isRequired {
			return fmt.Errorf("Unexpected mode for required mount point: %s", mode)
		}
		return nil
	}

	if err := os.MkdirAll(mountPoint.dst, os.ModePerm); err != nil {
		return err
	}
	if err := unix.Mount(mountPoint.src, mountPoint.dst, "bind", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return err
	}
	if mountPoint.isReadOnly {
		return unix.Mount("", mountPoint.dst, "", unix.MS_REMOUNT|unix.MS_BIND|unix.MS_RDONLY, "")
	}
	return nil
}

func (root *linuxRootFS) PivotRoot() error {

	if err := os.Chdir(root.target); err != nil {
		return err
	}

	oldRootTarget, err := createOldRootTarget()
	if err != nil {
		return err
	}
	if err := unix.PivotRoot(".", oldRootTarget); err != nil {
		return err
	}
	if err := removeOldRoot(oldRootTarget); err != nil {
		return err
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

func removeOldRoot(oldRootDir string) error {
	if err := unix.Unmount(oldRootDir, unix.MNT_DETACH); err != nil {
		return err
	}
	return os.Remove(oldRootDir)
}
