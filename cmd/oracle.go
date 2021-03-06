/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// oracleCmd represents the oracle command
var oracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "oracle sub-module under db module will Work with Oracle database layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	dbCmd.AddCommand(oracleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oracleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oracleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
