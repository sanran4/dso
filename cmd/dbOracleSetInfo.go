package cmd

import "github.com/spf13/cobra"

var dbOracleSetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		dbOrclSetCmd.Help()
		db_oracle_set_ex()
	},
}

func init() {
	dbOrclSetCmd.AddCommand(dbOracleSetInfoCmd)
}
