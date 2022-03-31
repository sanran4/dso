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
	Short: "This sub-module vcenter will work with VMware vCenter layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso virtual vcenter -h")
		fmt.Println("dso virtual vcenter --help")
		fmt.Println("dso help virtual vcenter")
	},
}

func init() {
	virtualCmd.AddCommand(vcenterCmd)

}
