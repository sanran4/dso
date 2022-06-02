/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/da0x/golang/olog"
	"github.com/jszwec/csvutil"
	"github.com/sanran4/dso/util"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbOrclReportCmd = &cobra.Command{
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
		outFormat, _ := cmd.Flags().GetString("out")
		srv, usr, pas, svc, prt := parseDbOrclRptFlags(cmd, args)

		// oracle://user:pass@server/service_name
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		//fmt.Println(connString)

		query1 := `
		select FILE_NAME, TABLESPACE_NAME, (BYTES/1048576) as SizeInMb from dba_data_files
		UNION
		select FILE_NAME, TABLESPACE_NAME, (BYTES/1048576) as SizeInMb from dba_temp_files
		UNION
		select name AS FILE_NAME, 'CONTROL' AS TABLESPACE_NAME, ((block_size * file_size_blks)/1048576) AS SizeInMb from v$controlfile WHERE IS_RECOVERY_DEST_FILE = 'NO'
		`
		fmt.Println("Oracle database data files details")
		out1 := getOrclDataFile(connString, query1)
		if outFormat == "table" {
			olog.Print(out1)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out1))
		} else if outFormat == "csv" {
			of1 := util.GetFilenameDate("oracleFilesReport", "csv")
			b1, err := csvutil.Marshal(out1)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of1, string(b1))
		}

		// select l.GROUP#, l.THREAD#, f.MEMBER, l.BYTES from v\$logfile f, v\$log l where f.group#=l.group#
		query2 := `
		select l.GROUP#, l.THREAD#, f.MEMBER, (l.BYTES/1048576) as SizeInMb from v$logfile f, v$log l where f.group# = l.group#
		`
		fmt.Println("Oracle database log files details")
		//getOrclFileDetails(connString, query2)
		out2 := getOrclLogFile(connString, query2)
		if outFormat == "table" {
			olog.Print(out2)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out2))
		} else if outFormat == "csv" {
			of2 := util.GetFilenameDate("oracleLogFilesReport", "csv")
			b2, err := csvutil.Marshal(out2)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of2, string(b2))
		}

		query3 := `
		select g.NAME, d.path,  d.total_mb, d.free_MB FROM v$asm_disk d, v$asm_diskgroup g where g.GROUP_NUMBER = d.GROUP_NUMBER
		`
		fmt.Println("Oracle ASM disk details")
		out3 := getOrclAsmDisks(connString, query3)
		if outFormat == "table" {
			olog.Print(out3)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out3))
		} else if outFormat == "csv" {
			of3 := util.GetFilenameDate("oracleLAsmDiskReport", "csv")
			b3, err := csvutil.Marshal(out3)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of3, string(b3))
		}

		query4 := `
		select NAME, DESCRIPTION, VALUE from v$parameter where name IN ('instance_type', 'instance_mode', 'compatible', 'compatible', 'service_names', 'db_name', 'processes', 'sessions', 'cpu_count', 'sga_min_size', 'sga_max_size', 'sga_target', 'db_block_size', 'memoptimize_pool_size', 'hash_area_size', 'result_cache_max_size', 'object_cache_optimal_size', 'sort_area_size', 'use_large_pages', 'log_buffer', 'background_dump_dest', 'user_dump_dest', 'core_dump_dest', 'audit_file_dest', 'optimizer_features_enable', 'parallel_degree_limit', 'enable_automatic_maintenance_pdb') ORDER BY NAME ASC
		`
		fmt.Println("Oracle database parameters detail")
		out4 := getOrclDbParameters(connString, query4)
		if outFormat == "table" {
			olog.Print(out4)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out4))
		} else if outFormat == "csv" {
			of4 := util.GetFilenameDate("oracleDbParameterReport", "csv")
			b4, err := csvutil.Marshal(out4)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of4, string(b4))
		}

	},
}

func init() {
	oracleCmd.AddCommand(dbOrclReportCmd)

	dbOrclReportCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclReportCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclReportCmd.Flags().StringP("instance", "I", "", "oracle db server name/IP address")
	dbOrclReportCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclReportCmd.Flags().String("svc", "", "oracle service name")
	dbOrclReportCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")

	//dbOrclRptCmd.MarkFlagRequired("server")
	//dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclReportCmd.MarkFlagRequired("svc")
}

func parseDbOrclRptFlags(cmd *cobra.Command, args []string) (srv, usr, pas, svc string, prt int) {
	server, ok := os.LookupEnv("ORCL_DB_HOST")
	if !ok {
		server, _ = cmd.Flags().GetString("instance")
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
	//reportType, _ := cmd.Flags().GetString("type")
	//outputType, _ := cmd.Flags().GetString("output")

	return server, user, pass, oraSvc, port
}

func dieOnError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

type orclDataFile struct {
	//Id   int64  `db:"name:visit_id"`
	FileName   string `db:"name:FILE_NAME"`
	TableSpace string `db:"name:TABLESPACE_NAME"`
	SizeInMB   string `db:"name:SizeInMb"`
	//Date time.Time	`db:"name:visit_date"`
}
type orclLogFile struct {
	GroupNo  int64  `db:"name:GROUP#"`
	ThreadNo string `db:"name:THREAD#"`
	Member   string `db:"name:MEMBER"`
	SizeInMB string `db:"name:SizeInMb"`
}

func getOrclLogFile(connStr, query string) []orclLogFile {
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
	//olog.Print(solf)
	return solf
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

func getOrclDbParameters(connStr, query string) []orclDbParameters {
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
	//olog.Print(sodf)
	return sodf
}

type orclAsmDisks struct {
	AsmGroupName  string `db:"name:NAME"`
	AsmDiskPath   string `db:"name:PATH"`
	TotalSizeInMB string `db:"name:Total_MB"`
	FreeSizeInMB  string `db:"name:Free_MB"`
}

func getOrclAsmDisks(connStr, query string) []orclAsmDisks {
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
	//olog.Print(sodf)
	return sodf
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
