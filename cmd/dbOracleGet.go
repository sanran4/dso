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

var dbOrclGetCmd = &cobra.Command{
	Use:   "get",
	Short: "This get command will Work for fetching best practice settings from oracle database layer ",
	Run: func(cmd *cobra.Command, args []string) {
		srv, usr, pas, svc, prt := parseDbOrclGetFlags(cmd, args)
		outFormat, _ := cmd.Flags().GetString("out")
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		//RestartOracleDatabase(connString)
		out1 := orclGetBPS(connString)
		fmt.Println("Oracle DB Best practice configuration:")
		if outFormat == "table" {
			olog.Print(out1)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out1))
		} else if outFormat == "csv" {
			of1 := util.GetFilenameDate("oracleBestPractice", "csv")
			b1, err := csvutil.Marshal(out1)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of1, string(b1))
		}

	},
}

func init() {
	oracleCmd.AddCommand(dbOrclGetCmd)

	dbOrclGetCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclGetCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclGetCmd.Flags().StringP("instance", "I", "", "oracle db server instance name/IP address")
	dbOrclGetCmd.Flags().Int("port", 1521, "oracle db port")
	dbOrclGetCmd.Flags().String("svc", "", "oracle service name")
	dbOrclGetCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")

	//dbOrclRptCmd.MarkFlagRequired("server")
	//dbOrclRptCmd.MarkFlagRequired("user")
	dbOrclGetCmd.MarkFlagRequired("svc")
}

func parseDbOrclGetFlags(cmd *cobra.Command, args []string) (srv, usr, pas, svc string, prt int) {
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
	//oraAttribute, _ = cmd.Flags().GetString("attr")

	return server, user, pass, oraSvc, port
}

type orclSgaOptimal struct {
	OptimalValue string `db:"name:OptimalValue"`
}

func orclGetSGAOptimal(connStr string) string {
	DB, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)
	defer DB.Close()

	query := `
	SELECT CAST((max(value)/1048576 )*0.8/64 AS integer)*64 AS OptimalValue
	FROM dba_hist_osstat WHERE stat_name = 'PHYSICAL_MEMORY_BYTES'
	`
	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	//rows, err := stmt.Query(nil)
	rows, err := stmt.Query_(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	var odf orclSgaOptimal
	for rows.Next_() {
		err = rows.Scan(&odf)
		dieOnError("Can't scan", err)
	}

	return odf.OptimalValue
}

type orclGetBPSData struct {
	//Id   int64  `db:"name:visit_id"`
	ConfigName   string `db:"name:ConfigName"`
	ConfigValue  string `db:"name:ConfigValue"`
	OptimalValue string `db:"name:OptimalValue"`
	Diff         string `db:"name:Diff"`
	//Date time.Time	`db:"name:visit_date"`
}

func orclGetBPS(connStr string) []orclGetBPSData {
	DB, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = DB.Open()
	dieOnError("Can't open the connection:", err)
	defer DB.Close()

	query := `
	select 'sga_max_size_MB' AS ConfigName, CAST((value/1048576) AS varchar2(50)) AS ConfigValue from v$parameter where name = 'sga_max_size'
	UNION ALL
	select 'sga_target_MB' AS ConfigName, CAST((value/1048576) AS varchar2(50)) AS ConfigValue from v$parameter where name = 'sga_target'
	UNION All
	select x.ksppinm ConfigName,  y.ksppstvl ConfigValue
	from sys.x$ksppi x, sys.x$ksppcv y
	where 1=1 and x.inst_id = y.inst_id and x.indx = y.indx
	and x.ksppinm ='_high_priority_processes'
	UNION All
	select '#LogFiles' AS ConfigName, CAST(count(1) AS varchar2(50)) AS ConfigValue from gv$log
	UNION All
	select 'LogFiles_Size_MB' AS ConfigName, CAST((max(BYTES)/1024/1024) AS varchar2(50)) AS ConfigValue from gv$log 
	`
	stmt := go_ora.NewStmt(query, DB)
	defer stmt.Close()

	//rows, err := stmt.Query(nil)
	rows, err := stmt.Query_(nil)
	dieOnError("Can't query", err)
	defer rows.Close()

	var odf orclGetBPSData
	var sodf []orclGetBPSData
	sgaSize := orclGetSGAOptimal(connStr)
	for rows.Next_() {
		err = rows.Scan(&odf.ConfigName, &odf.ConfigValue)
		dieOnError("Can't scan", err)
		switch odf.ConfigName {
		case "sga_max_size_MB":
			odf.OptimalValue = sgaSize
		case "sga_target_MB":
			odf.OptimalValue = sgaSize
		case "_high_priority_processes":
			odf.OptimalValue = "LMS*|VKTM|LGWR"
		case "#LogFiles":
			odf.OptimalValue = "5"
		case "LogFiles_Size_MB":
			odf.OptimalValue = "8192"
		}

		sodf = append(sodf, odf)

		for k := range sodf {
			if sodf[k].ConfigValue != sodf[k].OptimalValue {
				sodf[k].Diff = "*"
			}
		}
	}
	//fmt.Println(odf)
	//olog.Print(sodf)
	return sodf
}
