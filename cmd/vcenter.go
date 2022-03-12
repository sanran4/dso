/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// vcenterCmd represents the vcenter command
var vcenterCmd = &cobra.Command{
	Use:   "vcenter",
	Short: "Work with VMware vCenter layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso virtual vcenter -h")
		fmt.Println("dso virtual vcenter --help")
		fmt.Println("dso help virtual vcenter")
	},
}

func init() {
	virtualCmd.AddCommand(vcenterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vcenterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vcenterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
