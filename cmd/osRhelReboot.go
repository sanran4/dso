package cmd

import "github.com/spf13/cobra"

var osRhelRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot operating system remotely",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rhelCmd.AddCommand(osRhelRebootCmd)
}
