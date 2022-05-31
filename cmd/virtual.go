/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// virtualCmd represents the virtual command
var virtualCmd = &cobra.Command{
	Use:   "virtual",
	Short: "This virtual module will work with virtualization layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(virtualCmd)

}
