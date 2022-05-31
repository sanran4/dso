/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// hypervCmd represents the hyperv command
var hypervCmd = &cobra.Command{
	Use:   "hyperv",
	Short: "This hyperv sub-module will Work with Microsoft Hyper-V layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	virtualCmd.AddCommand(hypervCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hypervCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hypervCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
