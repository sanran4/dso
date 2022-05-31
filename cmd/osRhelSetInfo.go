package cmd

import "github.com/spf13/cobra"

var osRhelSetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		osRhelSetCmd.Help()
		os_rhel_set_ex()
	},
}

func init() {
	osRhelSetCmd.AddCommand(osRhelSetInfoCmd)
}
