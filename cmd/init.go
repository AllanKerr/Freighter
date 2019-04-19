package cmd

import (
	"sync"

	"github.com/allankerr/freighter/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {

		log.InitChildLogger()

		log.WithField("testField", 5).Info("info test")
		log.WithField("args", args).Trace("trace test")
		log.Debug("debug test")
		log.Warn("warn test")
		log.Error("error test")

		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	},
}
