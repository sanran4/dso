/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Get details and report about server layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand. please see help using below command:")
		fmt.Println("dso server -h")
		fmt.Println("dso server --help")
		fmt.Println("dso help server")
	},
}

func init() {
	serverCmd.AddCommand(serverSetCmd)

}
