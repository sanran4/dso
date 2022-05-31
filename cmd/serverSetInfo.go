package cmd

import "github.com/spf13/cobra"

var serverSetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		serverSetCmd.Help()
		server_set_ex()
	},
}

func init() {
	serverSetCmd.AddCommand(serverSetInfoCmd)
}
