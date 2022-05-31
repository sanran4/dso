package cmd

import "github.com/spf13/cobra"

var osRhelGetInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples of this command",
	Run: func(cmd *cobra.Command, args []string) {
		osRhelGetCmd.Help()
		os_rhel_get_ex()
	},
}

func init() {
	osRhelGetCmd.AddCommand(osRhelGetInfoCmd)
}
