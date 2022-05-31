package cmd

import "github.com/spf13/cobra"

var dbSqlGetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbSqlGetCmd.Help()
		db_sql_get_ex()
	},
}

func init() {
	dbSqlGetCmd.AddCommand(dbSqlGetInfoCmd)
}
