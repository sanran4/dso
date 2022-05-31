package cmd

import "github.com/spf13/cobra"

var dbSqlSetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbSqlSetCmd.Help()
		db_sql_set_ex()
	},
}

func init() {
	dbSqlSetCmd.AddCommand(dbSqlSetInfoCmd)
}
