package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sanran4/dso/util"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/spf13/cobra"
)

var oraAttribute string = ""

var dbOrclSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command will Work for changing oracle database layer settings",
	Run: func(cmd *cobra.Command, args []string) {
		srv, usr, pas, svc, prt := parseDbOrclSetFlags(cmd, args)
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		//RestartOracleDatabase(connString)
		if oraAttribute != "" {
			attr, val := parseOracleAttribute(oraAttribute)
			//if setParameter == ""
			if attr == "_high_priority_processes" {
				setHighPriorityProcess(connString, val)
			}
			if attr == "sga_target" || attr == "sga_max" {
				if val == ""{
					
				}else {
					setSgaValue(connString,attr,val)
				}
			}
		}
	},
}

func init() {
	oracleCmd.AddCommand(dbOrclSetCmd)

	dbOrclSetCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclSetCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclSetCmd.Flags().StringP("server", "S", "", "oracle db server name/IP address")
	dbOrclSetCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclSetCmd.Flags().String("svc", "", "oracle service name")
	dbOrclSetCmd.Flags().StringP("attr", "A", "", "Set individual attributes for oracle database")

	//dbOrclRptCmd.MarkFlagRequired("server")
	//dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclSetCmd.MarkFlagRequired("svc")
}

func parseOracleAttribute(str string) (attr, val string) {
	tmp := strings.Split(str, "=")
	attr = tmp[0]
	val = tmp[1]
	return
}

func parseDbOrclSetFlags(cmd *cobra.Command, args []string) (srv, usr, pas, svc string, prt int) {
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
	oraAttribute, _ = cmd.Flags().GetString("attr")

	return server, user, pass, oraSvc, port
}

func setSgaValue (connStr, attr, value string){

}

func setHighPriorityProcess(connStr, value string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()
	query1 := "alter system SET \"_high_priority_processes\" = '" + value + "' scope=spfile"
	stmt := go_ora.NewStmt(query1, conn)
	defer stmt.Close()
	result, err := stmt.Exec(nil)
	dieOnError("Can't query", err)
	fmt.Println(result)
}

func RestartOracleDatabase(connStr string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()

	//qry1 := "SHUTDOWN IMMEDIATE"
	//qry1 := "SHOW PARAMETER instance_name"
	//qry2 := `STARTUP`
	//stmt := go_ora.NewStmt(qry1, DB)
	//defer stmt.Close()
	stmt := go_ora.NewStmt("SHOW PARAMETER instance_name", conn)
	defer stmt.Close()
	fmt.Println(stmt)
	result, err := stmt.Exec(nil)
	//rows, err := stmt.Query_(nil)
	dieOnError("Can't query", err)
	//defer rows.Close()
	fmt.Println(result)

	//query1 = `
	//alter system SET "_high_priority_processes" = 'LMS*|VKTM|LGWR' scope=spfile
	//`
}
