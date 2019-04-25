package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/allankerr/freighter/controller"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		c := controller.NewController()
		if err := c.Run(); err != nil {
			log.WithError(err).Error("Failed to run container")
		} else {
			log.Info("Succeeded")
		}
	},
}
