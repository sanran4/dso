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
	Short: "This sub-module esxi Work with VMware ESXi layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso virtual esxi -h")
		fmt.Println("dso virtual esxi --help")
		fmt.Println("dso help virtual esxi")
	},
}

func init() {
	virtualCmd.AddCommand(esxiCmd)
}
