/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pstoreCmd represents the pstore command
var pstoreCmd = &cobra.Command{
	Use:   "pstore",
	Short: "Work with DellEMC PowerStore storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso storage pstore -h")
		fmt.Println("dso storage pstore --help")
		fmt.Println("dso help storage pstore")
	},
}

func init() {
	storageCmd.AddCommand(pstoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pstoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pstoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
