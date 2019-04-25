package tty

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type linuxTerminalMaster struct {
	Terminal
	master    *os.File
	slavePath string
}

func NewTerminalMaster() (TerminalMaster, error) {
	master, err := os.OpenFile("/dev/ptmx", unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	return NewTerminalMasterFile(master)

}

func NewTerminalMasterFile(master *os.File) (TerminalMaster, error) {
	slavePath, err := ptsname(master.Fd())
	if err != nil {
		return nil, err
	}
	terminal, err := NewTerminal(master)
	if err != nil {
		return nil, err
	}
	return &linuxTerminalMaster{
		Terminal:  terminal,
		master:    master,
		slavePath: slavePath,
	}, nil
}

func (t *linuxTerminalMaster) GetSlavePath() string {
	return t.slavePath
}

func (t *linuxTerminalMaster) UnlockSlave() error {
	return unlockpt(t.master.Fd())
}

func (t *linuxTerminalMaster) Listen() {
	go io.Copy(t.master, os.Stdin)
	go io.Copy(os.Stdout, t.master)
}

func unlockpt(masterFd uintptr) error {
	var unlock int32
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, masterFd, unix.TIOCSPTLCK, uintptr(unsafe.Pointer(&unlock))); err != 0 {
		return err
	}
	return nil
}

func ptsname(masterFd uintptr) (string, error) {
	var ptsno int32
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, masterFd, unix.TIOCGPTN, uintptr(unsafe.Pointer(&ptsno))); err != 0 {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", ptsno), nil
}
