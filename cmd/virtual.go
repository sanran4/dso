/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// virtualCmd represents the virtual command
var virtualCmd = &cobra.Command{
	Use:   "virtual",
	Short: "Work with virtualization layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso virtual -h")
		fmt.Println("dso virtual --help")
		fmt.Println("dso help virtual")
	},
}

func init() {
	rootCmd.AddCommand(virtualCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// virtualCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// virtualCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
