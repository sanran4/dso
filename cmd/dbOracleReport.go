/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sanran4/dso/util"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbOrclRptCmd = &cobra.Command{
	Use:   "report",
	Short: "Work with oracle database layer of the solution",
	//Usage: "dso db [command]",
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetInt("port")
		oraSvc, _ := cmd.Flags().GetString("svc")
		//database, _ := cmd.Flags().GetInt("database")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")
		var err error
		if pass == "" {
			pass, err = util.GetPasswd()
			if err != nil {
				log.Printf("error getting password %v", err)
			}
		}

		// oracle://user:pass@server/service_name
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", user, pass, server, port, oraSvc)
		fmt.Println(connString)
		// select l.GROUP#, l.THREAD#, f.MEMBER, l.BYTES from v\$logfile f, v\$log l where f.group#=l.group#
		cmd1 := `
		select l.GROUP#, l.THREAD#, f.MEMBER, l.BYTES from v\$logfile f, v\$log l where f.group#=l.group#
		`
		getOrclFileDetails(connString, cmd1)

		// select FILE_NAME, TABLESPACE_NAME, BYTES from dba_data_files;​
		cmd2 := `
		select FILE_NAME, TABLESPACE_NAME, BYTES from dba_data_files;
		`
		getOrclDataFile(connString, cmd2)
	},
}

func init() {
	oracleCmd.AddCommand(dbOrclRptCmd)

	dbOrclRptCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclRptCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclRptCmd.Flags().StringP("server", "S", "", "oracle db server name/IP address")
	dbOrclRptCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclRptCmd.Flags().String("svc", "", "oracle service name")

	dbOrclRptCmd.MarkFlagRequired("server")
	dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclRptCmd.MarkFlagRequired("svc")
}

func dieOnError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

type orclDataFile struct {
	//Id   int64  `db:"name:visit_id"`
	FileName    string `db:"name:FILE_NAME"`
	TableSpace  string `db:"name:TABLESPACE_NAME"`
	SizeInBytes string `db:"name:BYTES"`
	//Date time.Time	`db:"name:visit_date"`
}

func getOrclDataFile(connStr, query string) {
	DB, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)
	defer DB.Close()

	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	//rows, err := stmt.Query(nil)
	rows, err := stmt.Query_(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	var odf orclDataFile
	for rows.Next_() {
		err = rows.Scan(&odf)
		dieOnError("Can't scan", err)
		//fmt.Println("ID: ", Id, "\tName: ", vi.Name, "\tval: ", vi.Val, "\tDate: ", Date)
	}
	fmt.Println(odf)
}

func getOrclFileDetails(connStr, query string) {
	DB, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)
	defer DB.Close()

	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	rows, err := stmt.Query(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	columns := rows.Columns()
	values := make([]driver.Value, len(columns))
	Header(columns)

	for {
		err = rows.Next(values)
		if err != nil {
			break
		}
		Record(columns, values)
	}
	if err != io.EOF {
		dieOnError("Can't Next", err)
	}
}

func Header(columns []string) {

}

func Record(columns []string, values []driver.Value) {
	for i, c := range values {
		fmt.Printf("%-25s: %v\n", columns[i], c)
	}
	fmt.Println()
}
