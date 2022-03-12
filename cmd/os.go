/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// osCmd represents the os command
var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Work with operating system layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso os -h")
		fmt.Println("dso os --help")
		fmt.Println("dso help os")
	},
}

func init() {
	rootCmd.AddCommand(osCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// osCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// osCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
