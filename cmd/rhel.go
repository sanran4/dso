/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rhelCmd represents the rhel command
var rhelCmd = &cobra.Command{
	Use:   "rhel",
	Short: "This rhel module will Work with RHEL (Redhat Enterprise Linux) layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso os rhel -h")
		fmt.Println("dso os rhel --help")
		fmt.Println("dso help os rhel")
	},
}

func init() {
	osCmd.AddCommand(rhelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rhelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rhelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
