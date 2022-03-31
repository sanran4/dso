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
	Short: "This virtual module will work with virtualization layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing sub-module. please see help using below command:")
		fmt.Println("dso virtual -h")
		fmt.Println("dso virtual --help")
		fmt.Println("dso help virtual")
	},
}

func init() {
	rootCmd.AddCommand(virtualCmd)

}
