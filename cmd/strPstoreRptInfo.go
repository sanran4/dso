package cmd

import "github.com/spf13/cobra"

var pstoreRptInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		pstoreRptCmd.Help()
		storage_pstore_report_ex()
	},
}

func init() {
	pstoreRptCmd.AddCommand(pstoreRptInfoCmd)
}
