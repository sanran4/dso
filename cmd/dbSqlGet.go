package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/da0x/golang/olog"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jszwec/csvutil"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var dbSqlGetCmd = &cobra.Command{
	Use:   "get",
	Short: "This get command will validate best practice settings at SQL Server layer",
	Run: func(cmd *cobra.Command, args []string) {

		server, ok := os.LookupEnv("SQL_DB_HOST")
		if !ok {
			server, _ = cmd.Flags().GetString("instance")
		}
		user, ok := os.LookupEnv("SQL_DB_USER")
		if !ok {
			user, _ = cmd.Flags().GetString("user")
		}
		//server, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetInt("port")
		//database, _ := cmd.Flags().GetInt("database")
		//user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")
		var err error
		if pass == "" {
			pass, err = util.GetPasswd()
			if err != nil {
				log.Printf("error getting password %v", err)
			}
		}
		outFormat, _ := cmd.Flags().GetString("out")
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, pass, port)

		var query2 string = `
	if object_id('tempdb..#QueryOptimzer','u') is not null
		drop table #QueryOptimzer
	create table #QueryOptimzer (ConfigName sysname, ConfigValue varchar(16))
	
	exec sp_MSforeachdb 'Use ?
	INSERT INTO #QueryOptimzer
	select ''?'' +''_QUERY_OPTIMIZER_HOTFIXES'' as configName,convert(varchar(16),value) as ConfigValue 
	from sys.database_scoped_configurations where name = ''QUERY_OPTIMIZER_HOTFIXES''
	'
	select ConfigName, ConfigValue, OptimalValue,  
	CASE WHEN convert(bigint,ConfigValue) = convert(bigint,OptimalValue) THEN ''
	ELSE '*' END as Diff
	FROM (
	SELECT 'Dedicated Remote Admin Connection(DAC)'as ConfigName, value as ConfigValue, '1' as OptimalValue FROM sys.configurations WHERE name like '%remote admin connections%' 
	union all
	SELECT 'CPU_AffinityMask'as ConfigName, value as ConfigValue, (SELECT convert(varchar(50),(POWER(2,cpu_count)-1)) FROM sys.dm_os_sys_info) as OptimalValue
	FROM sys.configurations WHERE name like '%affinity mask%' 
	union all
	SELECT 'Max_server_memory'as ConfigName, value as ConfigValue, (SELECT convert(varchar(50),(total_physical_memory_kb/(1024))) from sys.dm_os_sys_memory) as OptimalValue 
	FROM sys.configurations WHERE name like '%max server memory%' 
	union all
	SELECT 'Min_server_memory'as ConfigName, value as ConfigValue, (SELECT convert(varchar(50),(total_physical_memory_kb/(1024))) from sys.dm_os_sys_memory) as OptimalValue
	FROM sys.configurations WHERE name like '%min server memory%'
	) tab1
	Union All
	select ConfigName, ConfigValue, OptimalValue, 
	CASE When ConfigValue = '1024 MB' Then ''
	ELSE '*' 
	END As Diff 
	FROM (
	SELECT db_name(database_id) + case type when 1 then '_LogFile_AutoGrowth' else '_DataFile_AutoGrowth' END as ConfigName,
	CASE is_percent_growth
	WHEN 1
	THEN CONVERT(varchar(16), growth) + '%'
	ELSE CONVERT(varchar(16), CONVERT(bigint, growth/128.0)) + ' MB'
	END as ConfigValue,
	'1024 MB' as OptimalValue
	FROM sys.master_files
	where database_id not in (db_id('master'),db_id('model'),db_id('msdb'),db_id('tempdb'))
	) tab2
	union 
	select ConfigName,ConfigValue,'1' as OptimalValue, 
	CASE When ConfigValue = '0' Then '*'
	ELSE '' 
	END As Diff 
	from #QueryOptimzer where ConfigName not in (
	'master_QUERY_OPTIMIZER_HOTFIXES',
	'tempdb_QUERY_OPTIMIZER_HOTFIXES',
	'model_QUERY_OPTIMIZER_HOTFIXES',
	'msdb_QUERY_OPTIMIZER_HOTFIXES')
	
		`
		fmt.Println("SQL Server Best practice configuration:")
		out2 := getSQLBPS(connString, query2)
		if outFormat == "table" {
			olog.Print(out2)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out2))
		} else if outFormat == "csv" {
			of2 := util.GetFilenameDate("sqlServerBPSConfig", "csv")
			b2, err := csvutil.Marshal(out2)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(of2, string(b2))
		}

	},
}

func init() {
	sqlCmd.AddCommand(dbSqlGetCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	dbSqlGetCmd.Flags().StringP("user", "U", "", "Username to connect to SQL Server instance")
	dbSqlGetCmd.Flags().StringP("pass", "P", "", "Password to connect to SQL Server instance")
	dbSqlGetCmd.Flags().StringP("instance", "I", "", "SQL Server instance name/IP address")
	dbSqlGetCmd.Flags().Int("port", 1433, "SQL Server instance port")
	dbSqlGetCmd.Flags().String("db", "", "SQL Server database name")
	dbSqlGetCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//dbSqlGetCmd.MarkFlagRequired("instance")
	//dbSqlGetCmd.MarkFlagRequired("user")
	//dbSqlReportCmd.MarkFlagRequired("pass")
}

type sqlBpsConfig struct {
	ConfigName   string `json:"ConfigName"`
	ConfigValue  string `json:"ConfigValue"`
	OptimalValue string `json:"OptimalValue"`
	Diff         string `json:"Diff"`
}

func getSQLBPS(connStr, query string) []sqlBpsConfig {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	totalRows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer totalRows.Close()

	mc1 := []sqlBpsConfig{}
	for totalRows.Next() {
		c1 := sqlBpsConfig{}
		err = totalRows.Scan(&c1.ConfigName, &c1.ConfigValue, &c1.OptimalValue, &c1.Diff)
		if err != nil {
			panic(err)
		}

		mc1 = append(mc1, c1)
	}
	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
	//olog.Print(mc1)
	return mc1
}
