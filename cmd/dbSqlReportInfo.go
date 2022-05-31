package cmd

import "github.com/spf13/cobra"

var dbSqlReportInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbSqlReportCmd.Help()
		db_sql_report_ex()
	},
}

func init() {
	dbSqlReportCmd.AddCommand(dbSqlReportInfoCmd)
}
