/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "This db module will Work with database layer of the solution",
	//Usage: "dso db [command]",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing sub-modules. please see help using below command:")
		fmt.Println("dso db -h")
		fmt.Println("dso db --help")
		fmt.Println("dso help db")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
