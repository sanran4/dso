/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// vcenterCmd represents the vcenter command
var vcenterCmd = &cobra.Command{
	Use:   "vcenter",
	Short: "This sub-module vcenter will work with VMware vCenter layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	virtualCmd.AddCommand(vcenterCmd)

}
