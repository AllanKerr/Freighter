package cmd

import (
	"fmt"

	"github.com/allankerr/freighter/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		freighter, err := cli.NewDefaultConfiguration()
		if err != nil {
			fmt.Println("Unable to initialize Freighter:", err)
			return
		}
		if err := freighter.Create(args[0], args[1]); err != nil {
			fmt.Println("Unable to create container:", err)
		}
	},
}
