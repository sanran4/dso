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
	Short: "Work with DellEMC PowerMax storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso storage pmax -h")
		fmt.Println("dso storage pmax --help")
		fmt.Println("dso help storage pmax")
	},
}

func init() {
	storageCmd.AddCommand(pmaxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pmaxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pmaxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
