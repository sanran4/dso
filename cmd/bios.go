package cmd

import "github.com/spf13/cobra"

var serverBiosCmd = &cobra.Command{
	Use:   "bios",
	Short: "bios sub-module under server module provides commands to work with BIOS settings of the physical server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	serverCmd.AddCommand(serverBiosCmd)
}
