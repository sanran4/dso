/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "This module storage will Work with DellEMC storage layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing sub-module. please see help using below command:")
		fmt.Println("dso storage -h")
		fmt.Println("dso storage --help")
		fmt.Println("dso help storage")
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)

}
