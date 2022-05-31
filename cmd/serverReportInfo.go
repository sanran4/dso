package cmd

import "github.com/spf13/cobra"

var serverReportInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		serverReportCmd.Help()
		server_report_ex()
	},
}

func init() {
	serverReportCmd.AddCommand(serverReportInfoCmd)
}
