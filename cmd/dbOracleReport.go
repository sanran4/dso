/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/da0x/golang/olog"
	"github.com/sanran4/dso/util"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbOrclRptCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will Work for reporting oracle database layer settings",
	//Usage: "dso db [command]",
	Run: func(cmd *cobra.Command, args []string) {
		/*
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
		*/

		srv, usr, pas, svc, prt, rptType, outType := parseDbOrclRptFlags(cmd, args)

		// oracle://user:pass@server/service_name
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		//fmt.Println(connString)

		// select FILE_NAME, TABLESPACE_NAME, BYTES from dba_data_files;​
		if rptType == "" || rptType == "dbfile" {
			q1 := `
			select FILE_NAME, TABLESPACE_NAME, BYTES from dba_data_files
			`
			fmt.Println("Oracle database data files details")
			out := getOrclDataFile(connString, q1)
			if outType == "table" {
				olog.Print(out)
			} else if outType == "pdf" {
				CreatePDF("Oracle database data files details", out)
			}
		}

		query1 := `
		select FILE_NAME, TABLESPACE_NAME, BYTES from dba_data_files
		UNION
		select FILE_NAME, TABLESPACE_NAME, BYTES from dba_temp_files
		UNION
		select name AS FILE_NAME, 'CONTROL' AS TABLESPACE_NAME, (block_size * file_size_blks) AS BYTES from v$controlfile WHERE IS_RECOVERY_DEST_FILE = 'NO'
		`
		fmt.Println("Oracle database data files details")
		getOrclDataFile(connString, query1)

		// select l.GROUP#, l.THREAD#, f.MEMBER, l.BYTES from v\$logfile f, v\$log l where f.group#=l.group#
		query2 := `
		select l.GROUP#, l.THREAD#, f.MEMBER, l.BYTES from v$logfile f, v$log l where f.group# = l.group#
		`
		fmt.Println("Oracle database log files details")
		//getOrclFileDetails(connString, query2)
		getOrclLogFile(connString, query2)

		query3 := `
		select g.NAME, d.path,  d.total_mb FROM v$asm_disk d, v$asm_diskgroup g where g.GROUP_NUMBER = d.GROUP_NUMBER
		`
		fmt.Println("Oracle ASM disk details")
		getOrclAsmDisks(connString, query3)

		query4 := `
		select NAME, DESCRIPTION, VALUE from v$parameter where name IN ('instance_type', 'instance_mode', 'compatible', 'compatible', 'service_names', 'db_name', 'processes', 'sessions', 'cpu_count', 'sga_min_size', 'sga_max_size', 'sga_target', 'db_block_size', 'memoptimize_pool_size', 'hash_area_size', 'result_cache_max_size', 'object_cache_optimal_size', 'sort_area_size', 'use_large_pages', 'log_buffer', 'background_dump_dest', 'user_dump_dest', 'core_dump_dest', 'audit_file_dest', 'optimizer_features_enable', 'parallel_degree_limit', 'enable_automatic_maintenance_pdb') ORDER BY NAME ASC
		`
		fmt.Println("Oracle database parameters detail")
		getOrclDbParameters(connString, query4)

	},
}

func init() {
	oracleCmd.AddCommand(dbOrclRptCmd)

	dbOrclRptCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclRptCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclRptCmd.Flags().StringP("server", "S", "", "oracle db server name/IP address")
	dbOrclRptCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclRptCmd.Flags().String("svc", "", "oracle service name")
	dbOrclRptCmd.Flags().StringP("type", "t", "", "Report type to be fetched")
	dbOrclRptCmd.Flags().StringP("output", "o", "table", "Output report data in table/pdf format")

	//dbOrclRptCmd.MarkFlagRequired("server")
	//dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclRptCmd.MarkFlagRequired("svc")
}

func parseDbOrclRptFlags(cmd *cobra.Command, args []string) (srv, usr, pas, svc string, prt int, rptTyp, outTyp string) {
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
	reportType, _ := cmd.Flags().GetString("type")
	outputType, _ := cmd.Flags().GetString("output")

	return server, user, pass, oraSvc, port, reportType, outputType
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
type orclLogFile struct {
	GroupNo     int64  `db:"name:GROUP#"`
	ThreadNo    string `db:"name:THREAD#"`
	Member      string `db:"name:MEMBER"`
	SizeInBytes string `db:"name:BYTES"`
}

func getOrclLogFile(connStr, query string) {
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

	var olf orclLogFile
	var solf []orclLogFile
	for rows.Next_() {
		err = rows.Scan(&olf)
		dieOnError("Can't scan", err)
		solf = append(solf, olf)
	}
	olog.Print(solf)
}

func getOrclDataFile(connStr, query string) []orclDataFile {
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
	var sodf []orclDataFile
	for rows.Next_() {
		err = rows.Scan(&odf)
		dieOnError("Can't scan", err)
		sodf = append(sodf, odf)
		//fmt.Println("ID: ", Id, "\tName: ", vi.Name, "\tval: ", vi.Val, "\tDate: ", Date)
		//fmt.Println(odf)
	}
	//fmt.Println(odf)
	//olog.Print(sodf)
	return sodf
}

type orclDbParameters struct {
	ParameterName string `db:"name:NAME"`
	Discription   string `db:"name:DESCRIPTION"`
	CurrentValue  string `db:"name:VALUE"`
}

func getOrclDbParameters(connStr, query string) {
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

	var odf orclDbParameters
	var sodf []orclDbParameters
	for rows.Next_() {
		err = rows.Scan(&odf)
		dieOnError("Can't scan", err)
		sodf = append(sodf, odf)
		//fmt.Println("ID: ", Id, "\tName: ", vi.Name, "\tval: ", vi.Val, "\tDate: ", Date)
		//fmt.Println(odf)
	}
	//fmt.Println(odf)
	olog.Print(sodf)
}

type orclAsmDisks struct {
	AsmGroupName string `db:"name:NAME"`
	AsmDiskPath  string `db:"name:PATH"`
	SizeInMB     string `db:"name:Total_MB"`
}

func getOrclAsmDisks(connStr, query string) {
	DB, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)
	defer DB.Close()

	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	rows, err := stmt.Query_(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	var odf orclAsmDisks
	var sodf []orclAsmDisks
	for rows.Next_() {
		err = rows.Scan(&odf)
		dieOnError("Can't scan", err)
		sodf = append(sodf, odf)
	}
	//fmt.Println(odf)
	olog.Print(sodf)
}

func CreatePDF(reportHeading string, reportData []orclDataFile) {

}

/*
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
*/
