package cmd

import (
	"github.com/allankerr/freighter/process"
	"github.com/allankerr/freighter/spec"
	"github.com/allankerr/freighter/uts"

	"github.com/allankerr/freighter/fs"

	"github.com/allankerr/freighter/log"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {

		unix.Setsid()
		log.InitChildLogger()
		log.Debug("Child process started")

		f, err := unix.Open("/proc/self/fd/4", unix.O_RDWR, 0)
		if err != nil {
			log.Fatal("Failed to open pty slave")
		}
		for stdfd := 0; stdfd < 3; stdfd++ {
			unix.Dup2(f, stdfd)
		}

		config := spec.BaseConfig

		rootfs, err := fs.NewRootFS(config.Root)
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

		if err := uts.SetHostname(config.Hostname); err != nil {
			log.WithError(err).Fatal("Failed to set hostname")
		}

		proc := process.New(config.Process)
		if err := proc.Run(); err != nil {
			log.WithError(err).Fatal("Failed to run process")
		}
	},
}
