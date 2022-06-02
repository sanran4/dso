package cmd

import "github.com/spf13/cobra"

var serverBiosCmd = &cobra.Command{
	Use:   "bios",
	Short: "bios sub-module will work with BIOS settings of the physical server layer within solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	serverCmd.AddCommand(serverBiosCmd)
}
