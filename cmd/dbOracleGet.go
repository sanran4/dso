package cmd

import (
	"log"
	"os"

	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

var dbOrclGetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command will Work for changing oracle database layer settings",
	Run: func(cmd *cobra.Command, args []string) {
		//srv, usr, pas, svc, prt := parseDbOrclGetFlags(cmd, args)
		//connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		//RestartOracleDatabase(connString)

	},
}

func init() {
	oracleCmd.AddCommand(dbOrclGetCmd)

	dbOrclGetCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclGetCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclGetCmd.Flags().StringP("server", "S", "", "oracle db server name/IP address")
	dbOrclGetCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclGetCmd.Flags().String("svc", "", "oracle service name")
	//dbOrclGetCmd.Flags().StringP("attr", "A", "", "Set individual attributes for oracle database")

	//dbOrclRptCmd.MarkFlagRequired("server")
	//dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclGetCmd.MarkFlagRequired("svc")
}

func parseDbOrclGetFlags(cmd *cobra.Command, args []string) (srv, usr, pas, svc string, prt int) {
	server, ok := os.LookupEnv("ORCL_DB_HOST")
	if !ok {
		server, _ = cmd.Flags().GetString("server")
	}
	user, ok := os.LookupEnv("ORCL_DB_USER")
	if !ok {
		user, _ = cmd.Flags().GetString("user")
	}

	pass, _ := cmd.Flags().GetString("pass")
	var err error
	if pass == "" {
		pass, err = util.GetPasswd()
		if err != nil {
			log.Printf("error getting password %v", err)
		}
	}

	oraSvc, _ := cmd.Flags().GetString("svc")
	port, _ := cmd.Flags().GetInt("port")
	//oraAttribute, _ = cmd.Flags().GetString("attr")

	return server, user, pass, oraSvc, port
}
