/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// esxiCmd represents the esxi command
var esxiCmd = &cobra.Command{
	Use:   "esxi",
	Short: "This sub-module esxi Work with VMware ESXi layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	virtualCmd.AddCommand(esxiCmd)
}
