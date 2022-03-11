/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// biosCmd represents the bios command
var biosCmd = &cobra.Command{
	Use:   "bios",
	Short: "fetch report or change setting for server BIOS",
	Run: func(cmd *cobra.Command, args []string) {
		//parseBiosArgs(cmd, args)
	},
}

func init() {
	serverCmd.AddCommand(biosCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	biosCmd.PersistentFlags().BoolP("report", "r", true, "Fetch report for the server bios")
	biosCmd.PersistentFlags().StringP("setConfig", "s", "", "Password for the server iDRAC")

}
