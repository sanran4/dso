/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// pmaxCmd represents the pmax command
var pmaxCmd = &cobra.Command{
	Use:   "pmax",
	Short: "This pmax sub-module will Work with DellEMC PowerMax storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	storageCmd.AddCommand(pmaxCmd)
}
