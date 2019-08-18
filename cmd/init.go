package cmd

import (
	"github.com/allankerr/freighter/container"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		c := container.NewContainer()
		c.Initialize()
	},
}
