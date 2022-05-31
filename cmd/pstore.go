/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// pstoreCmd represents the pstore command
var pstoreCmd = &cobra.Command{
	Use:   "pstore",
	Short: "This pstore sub-module will Work with DellEMC PowerStore storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	storageCmd.AddCommand(pstoreCmd)

}
