/*
Copyright © 2022 Sanjeev Ranjan

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dso",
	Short: "A database solution optimizer tool/utility",
	Long: `dso is a database solution optimizer tool that can be used to:
	- report different configuration settings values at given layer within solution/s
	- validate best practice recomendations for a given layer within solution/s
	- apply best practice setting/s at individual layer within database soluton
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		chkVersion, _ := cmd.Flags().GetBool("version")
		if chkVersion {
			fmt.Println("DSO Version v1.0.0")
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dso.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "get current version of the dso")
}
