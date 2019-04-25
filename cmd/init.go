package cmd

import (
	"os"
	"syscall"

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

		log.Debug("Calling exec")
		err = syscall.Exec("/bin/bash", []string{"/bin/bash"}, os.Environ())
		if err != nil {
			log.WithError(err).Error("exec failed")
		}
	},
}
