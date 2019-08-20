package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/allankerr/freighter/ipc"
	"github.com/allankerr/freighter/log"
	"github.com/allankerr/freighter/spec"

	"github.com/allankerr/freighter/cli/state"
	"golang.org/x/sys/unix"
)

func createContainer(state *state.State, config *spec.Spec) error {

	log.Error(state.ID())
	fifo, err := createFIFO(state.DirectoryPath())
	if err != nil {
		return err
	}
	initP, initC, err := createInitPipe()
	if err != nil {
		return err
	}
	logC, err := createRemoteLogger(state.ID())
	if err != nil {
		return err
	}

	flags := 0
	for _, namespace := range config.Linux.Namespaces {
		flags |= namespace.GetType()
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.Env = []string{}
	cmd.ExtraFiles = []*os.File{}
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: uintptr(flags),
	}

	appendFile(cmd, fifo, ipc.FifoName)
	appendFile(cmd, initC, ipc.InitName)
	appendFile(cmd, logC, ipc.LogName)

	if err := cmd.Start(); err != nil {
		return err
	}

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	initP.Write(configData)
	initP.WriteString("\n")

	_, err = cmd.Process.Wait()
	return err
}

func appendFile(cmd *exec.Cmd, file *os.File, name string) {

	fd := uintptr(len(cmd.ExtraFiles)) + os.Stderr.Fd() + 1
	envVar := fmt.Sprintf("%s=%v", name, fd)
	cmd.Env = append(cmd.Env, envVar)
	cmd.ExtraFiles = append(cmd.ExtraFiles, file)
}

func createFIFO(directoryPath string) (*os.File, error) {

	path := path.Join(directoryPath, "fifo")
	if err := syscall.Mkfifo(path, uint32(os.ModePerm)); err != nil {
		return nil, err
	}
	return os.OpenFile(path, unix.O_PATH|unix.FD_CLOEXEC, os.ModePerm)
}

func createInitPipe() (*os.File, *os.File, error) {

	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, err
	}
	initC := os.NewFile(uintptr(fds[1]), "init-c")
	initP := os.NewFile(uintptr(fds[0]), "init-p")
	unix.CloseOnExec(fds[1])
	unix.CloseOnExec(fds[0])
	return initP, initC, nil
}

func createRemoteLogger(containerID string) (*os.File, error) {
	remoteLogger, err := log.NewRemoteLogger(containerID)
	if err != nil {
		return nil, err
	}
	remoteLogger.Listen()
	return remoteLogger.Child(), nil
}
