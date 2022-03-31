/*
Copyright © 2022 Sanjeev Ranjan <s_ranjan@dell.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// windowsCmd represents the windows command
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "this sub-module will work with Microsoft Windows Server layer of the solution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing command. please see help using below command:")
		fmt.Println("dso os windows -h")
		fmt.Println("dso os windows --help")
		fmt.Println("dso help os windows")
	},
}

func init() {
	osCmd.AddCommand(windowsCmd)

}
