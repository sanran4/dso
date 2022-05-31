package cmd

import "github.com/spf13/cobra"

var osRhelReportInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		osRhelReportCmd.Help()
		os_rhel_report_ex()
	},
}

func init() {
	osRhelReportCmd.AddCommand(osRhelReportInfoCmd)
}
