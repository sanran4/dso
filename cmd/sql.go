/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "This sub-module sql will Work with Microsoft SQL Server database layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso db sql -h")
		fmt.Println("dso db sql --help")
		fmt.Println("dso help db sql")
	},
}

func init() {
	dbCmd.AddCommand(sqlCmd)

}
