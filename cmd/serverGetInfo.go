package cmd

import "github.com/spf13/cobra"

var serverGetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		serverGetCmd.Help()
		server_get_ex()
	},
}

func init() {
	serverGetCmd.AddCommand(serverGetInfoCmd)
}
