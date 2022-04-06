/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pmaxCmd represents the pmax command
var pmaxCmd = &cobra.Command{
	Use:   "pmax",
	Short: "This pmax sub-module will Work with DellEMC PowerMax storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso storage pmax -h")
		fmt.Println("dso storage pmax --help")
		fmt.Println("dso help storage pmax")
	},
}

func init() {
	storageCmd.AddCommand(pmaxCmd)
}
