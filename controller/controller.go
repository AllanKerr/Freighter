package controller

import (
	"os"
	"os/exec"

	"github.com/allankerr/freighter/log"
	"github.com/allankerr/freighter/tty"
	"golang.org/x/sys/unix"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Run() error {

	remoteLogger := log.InitParentLogger()

	curPty, err := tty.NewTerminal(os.Stdin)
	if err != nil {
		return err
	}
	if err := curPty.SetRaw(); err != nil {
		return err
	}
	defer curPty.Reset()

	containerPty, err := tty.NewTerminalMaster()
	if err != nil {
		return err
	}
	if err := containerPty.UnlockSlave(); err != nil {
		return err
	}
	slave, err := os.OpenFile(containerPty.GetSlavePath(), unix.O_PATH, 0)
	if err != nil {
		return err
	}
	containerPty.Listen()

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWPID | unix.CLONE_NEWNS | unix.CLONE_NEWUTS,
	}

	cmd.ExtraFiles = []*os.File{remoteLogger.Child(), slave}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = cmd.Process.Wait()
	return err
}
