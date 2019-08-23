package container

import (
	"fmt"
	"os"

	"github.com/allankerr/freighter/cli/state"

	"github.com/allankerr/freighter/ipc"
	"github.com/allankerr/freighter/log"
	"github.com/allankerr/freighter/process"
	"github.com/allankerr/freighter/spec"
	"golang.org/x/sys/unix"
)

type Container interface {
	Initialize()
}

type containerLinux struct {
}

func NewContainer() Container {
	logFD, err := findFileDescriptor(ipc.LogName)
	if err == nil {
		file := os.NewFile(logFD, "log-c")
		log.InitChildLogger(file)
	}
	return &containerLinux{}
}

func (c *containerLinux) Initialize() {

	log.Debug("Child process started")

	unix.Setsid()

	/*f, err := unix.Open("/proc/self/fd/4", unix.O_RDWR, 0)
	if err != nil {
		log.Fatal("Failed to open pty slave")
	}
	for stdfd := 0; stdfd < 3; stdfd++ {
		unix.Dup2(f, stdfd)
	}*/

	initFD, err := findFileDescriptor(ipc.InitName)
	if err != nil {
		log.WithError(err).Fatal("Failed to find init file descriptor")
	}
	fifoFD, err := findFileDescriptor(ipc.FifoName)
	if err != nil {
		log.WithError(err).Fatal("Failed to find FIFO file descriptor")
	}

	init := os.NewFile(initFD, "init-c")
	initPipe := ipc.NewPipe(init)

	config := &spec.Spec{}
	msgType, err := initPipe.Receive(config)
	if err != nil {
		log.WithError(err).Fatal("Failed to read container configuration")
	}
	if msgType != ipc.MessageInitSpec {
		log.WithField("type", msgType).Fatal("Received unexpected message type")
	}
	log.WithField("config", config).Debug("Read config from init-c")

	rootfs, err := NewRootFS(config.Root)
	if err != nil {
		log.WithError(err).Fatal("Failed to create rootfs")
	}

	if err := rootfs.PrepareRoot(config.Linux.RootFSPropagation); err != nil {
		log.WithError(err).Fatal("Failed to prepare root")
	}
	if err := rootfs.CreateMounts(config.Mounts); err != nil {
		log.WithError(err).Fatal("Failed to create mount")
	}
	if err := rootfs.CreateDevices(config.Linux.Devices); err != nil {
		log.WithError(err).Fatal("Failed to create devices")
	}
	if err := rootfs.PivotRoot(); err != nil {
		log.WithError(err).Fatal("Failed to pivot root")
	}
	if err := rootfs.FinalizeRoot(); err != nil {
		log.WithError(err).Fatal("Failed to finalize root")
	}

	if err := setHostname(config.Hostname); err != nil {
		log.WithError(err).Fatal("Failed to set hostname")
	}

	log.Info("Container created...")

	statusChangePayload := &ipc.StatusChangePayload{
		Status: state.Created,
	}
	if err := initPipe.Send(ipc.MessageStatusChange, statusChangePayload); err != nil {
		log.WithError(err).Fatal("Failed to send status change message")
	}

	fifoPath := fmt.Sprintf("/proc/self/fd/%d", fifoFD)
	_, err = os.OpenFile(fifoPath, unix.O_WRONLY, 0)
	if err != nil {
		log.WithError(err).Fatal("Failed to open FIFO file descriptor")
	}

	proc := process.New(config.Process)
	if err := proc.Run(); err != nil {
		log.WithError(err).Fatal("Failed to run process")
	}
}
