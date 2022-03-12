/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// esxiCmd represents the esxi command
var esxiCmd = &cobra.Command{
	Use:   "esxi",
	Short: "Work with VMware ESXi layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso virtual esxi -h")
		fmt.Println("dso virtual esxi --help")
		fmt.Println("dso help virtual esxi")
	},
}

func init() {
	virtualCmd.AddCommand(esxiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// esxiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// esxiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
