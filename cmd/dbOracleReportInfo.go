package cmd

import "github.com/spf13/cobra"

var dbOracleRptInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbOrclRptCmd.Help()
		db_oracle_report_ex()
	},
}

func init() {
	dbOrclRptCmd.AddCommand(dbOracleRptInfoCmd)
}
