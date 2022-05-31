package cmd

import "github.com/spf13/cobra"

var dbOracleGetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbOrclGetCmd.Help()
		db_oracle_get_ex()
	},
}

func init() {
	dbOrclGetCmd.AddCommand(dbOracleGetInfoCmd)
}
