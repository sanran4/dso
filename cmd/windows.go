/*
Copyright © 2022 Sanjeev Ranjan <s_ranjan@dell.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// windowsCmd represents the windows command
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "this sub-module will work with Microsoft Windows Server layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	osCmd.AddCommand(windowsCmd)

}
