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
	Short: "This pstore sub-module will Work with DellEMC PowerStore storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso storage pstore -h")
		fmt.Println("dso storage pstore --help")
		fmt.Println("dso help storage pstore")
	},
}

func init() {
	storageCmd.AddCommand(pstoreCmd)

}
