/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// biosCmd represents the bios command
var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "fetch report or change setting for server BIOS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso bios -h")
		fmt.Println("dso bios --help")
		fmt.Println("dso help bios")
	},
}

func init() {
	serverCmd.AddCommand(biosCmd)

}
