package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "get help and examples for requested command",
	Run: func(cmd *cobra.Command, args []string) {
		fetchInfo(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func fetchInfo(cmd *cobra.Command, args []string) {
	arguments := args
	argLen := len(arguments)
	//var arg1, arg2, arg3 string
	for i := 0; i < argLen; i++ {
		compVar := i + 1
		if i == 0 {
			if compVar == argLen {
				switch arguments[i] {
				case "server":
					serverCmd.Help()
				case "os":
					osCmd.Help()
				case "db":
					dbCmd.Help()
				case "storage":
					storageCmd.Help()
				default:
					fmt.Println("invalid argument")
					return
				}
				return
			}
		}
		if i == 1 {
			if compVar == argLen {
				if arguments[0] == "server" {
					if arguments[1] == "bios" {
						serverBiosCmd.Help()
					} else if arguments[1] == "reboot" {
						serverRebootCmd.Help()
					} else {
						fmt.Println("invalid argument")
						return
					}
				}
				if arguments[0] == "os" {
					if arguments[1] == "rhel" {
						rhelCmd.Help()
					} else if arguments[1] == "reboot" {
						osRhelRebootCmd.Help()
					} else {
						fmt.Println("invalid argument")
						return
					}
				}
				if arguments[0] == "db" {
					if arguments[1] == "sql" {
						sqlCmd.Help()
					} else if arguments[1] == "oracle" {
						oracleCmd.Help()
					} else {
						fmt.Println("invalid argument")
						return
					}
				}
				if arguments[0] == "storage" {
					if arguments[1] == "pstore" {
						pstoreCmd.Help()
					} else {
						fmt.Println("invalid argument")
						return
					}
				}
				return
			}
		}
		if i == 2 {
			if compVar == argLen {
				if arguments[0] == "server" {
					if arguments[1] == "bios" {
						if arguments[2] == "get" {
							serverBiosGetCmd.Help()
							server_bios_get_ex()
						} else if arguments[2] == "set" {
							serverBiosSetCmd.Help()
							server_bios_set_ex()
						} else if arguments[2] == "report" {
							serverBiosReportCmd.Help()
							server_bios_report_ex()
						} else {
							fmt.Println("invalid argument")
							return
						}
					}
				} else if arguments[0] == "os" {
					if arguments[1] == "rhel" {
						if arguments[2] == "get" {
							osRhelGetCmd.Help()
							os_rhel_get_ex()
						} else if arguments[2] == "set" {
							osRhelSetCmd.Help()
							os_rhel_set_ex()
						} else if arguments[2] == "report" {
							osRhelReportCmd.Help()
							os_rhel_report_ex()
						} else {
							fmt.Println("invalid argument")
							return
						}
					}
				} else if arguments[0] == "db" {
					if arguments[1] == "sql" {
						if arguments[2] == "get" {
							dbSqlGetCmd.Help()
							db_sql_get_ex()
						} else if arguments[2] == "set" {
							dbSqlSetCmd.Help()
							db_sql_set_ex()
						} else if arguments[2] == "report" {
							dbSqlReportCmd.Help()
							db_sql_report_ex()
						} else {
							fmt.Println("invalid argument")
							return
						}
					}
					if arguments[1] == "oracle" {
						if arguments[2] == "get" {
							dbOrclGetCmd.Help()
							db_oracle_get_ex()
						} else if arguments[2] == "set" {
							dbOrclSetCmd.Help()
							db_oracle_set_ex()
						} else if arguments[2] == "report" {
							dbOrclReportCmd.Help()
							db_oracle_report_ex()
						} else {
							fmt.Println("invalid argument")
							return
						}
					}
				} else if arguments[0] == "storage" {
					if arguments[1] == "pstore" {
						if arguments[2] == "report" {
							pstoreCmd.Help()
							storage_pstore_report_ex()
						} else {
							fmt.Println("invalid argument")
							return
						}
					}
				} else {
					fmt.Println("invalid argument")
					return
				}
				return
			}

		}
		if i > 2 {
			fmt.Println("invalid argument")
			return
		}
	}
}
