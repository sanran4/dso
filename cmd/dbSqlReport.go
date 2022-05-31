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
var dbSqlReportCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will pull report from SQL Server",
	Long:  `This report command will pull report from SQL Server`,
	Run: func(cmd *cobra.Command, args []string) {

		server, ok := os.LookupEnv("SQL_DB_HOST")
		if !ok {
			server, _ = cmd.Flags().GetString("instance")
		}
		user, ok := os.LookupEnv("SQL_DB_USER")
		if !ok {
			user, _ = cmd.Flags().GetString("user")
		}
		//user, _ := cmd.Flags().GetString("user")
		port, _ := cmd.Flags().GetInt("port")
		//database, _ := cmd.Flags().GetInt("database")

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

		var querySqlServerProperty string = `
		select srvSetting, srvData from (
			SELECT  
			  SERVERPROPERTY('MachineName') AS ComputerName,
			  SERVERPROPERTY('ServerName') AS InstanceName,  
			  SERVERPROPERTY('ProductVersion') AS ProductVersion,  
			  SERVERPROPERTY('ProductLevel') AS ProductLevel,  
			  SERVERPROPERTY('Edition') AS Edition,
			  SERVERPROPERTY('InstanceDefaultDataPath') AS InstanceDefaultDataPath,
			  SERVERPROPERTY('InstanceDefaultLogPath') AS InstanceDefaultLogPath,
			  SERVERPROPERTY('InstanceDefaultBackupPath') AS InstanceDefaultBackupPath,
			  SERVERPROPERTY('Collation') AS Collation,
			  SERVERPROPERTY('IsClustered') AS IsClustered,
			  SERVERPROPERTY('IsHadrEnabled') AS IsHadrEnabled,
			  SERVERPROPERTY('IsPolyBaseInstalled') AS IsPolyBaseInstalled
			) as t1
			UNPIVOT(srvData FOR srvSetting IN (
				ComputerName, InstanceName, ProductVersion, ProductLevel, Edition, InstanceDefaultDataPath, InstanceDefaultLogPath, InstanceDefaultBackupPath, Collation, IsClustered, IsHadrEnabled, IsPolyBaseInstalled)
				) AS unp
		`
		fmt.Println("SQL Server Property:")
		output1 := getServerConfig(connString, querySqlServerProperty)
		if outFormat == "table" {
			olog.Print(output1)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(output1))
		} else if outFormat == "csv" {
			outputFile1 := util.GetFilenameDate("SqlServerPropertyReport", "csv")
			b1, err := csvutil.Marshal(output1)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile1, string(b1))
		}

		var queryResource string = `
		select setting,currentValue from (
		SELECT convert(varchar(50),CONVERT(varchar,sqlserver_start_time,20))sqlserver_start_time,convert(varchar(50),socket_count)socket_count,convert(varchar(50),cores_per_socket)cores_per_socket,
		convert(varchar(50),cpu_count)cpu_count, convert(varchar(50),softnuma_configuration_desc)softnuma_configuration_desc, convert(varchar(50),numa_node_count)numa_node_count, 
		convert(varchar(50),scheduler_total_count)scheduler_total_count, convert(varchar(50),scheduler_count)scheduler_count, convert(varchar(50),affinity_type_desc)affinity_type_desc,
		convert(varchar(50),process_physical_affinity)process_physical_affinity,convert(varchar(50),(physical_memory_kb/(1024*1024)))physical_memory_gb,convert(varchar(50),max_workers_count)max_workers_count   
		FROM sys.dm_os_sys_info
		) as t1
		UNPIVOT(currentValue FOR setting IN (
		sqlserver_start_time, socket_count, cores_per_socket, cpu_count,softnuma_configuration_desc,numa_node_count, scheduler_total_count,scheduler_count,affinity_type_desc,process_physical_affinity,
		physical_memory_gb,max_workers_count)
		) AS unp
		`
		fmt.Println("SQL Server Resource:")
		out2 := getDbSettingConfig(connString, queryResource)
		if outFormat == "table" {
			olog.Print(out2)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out2))
		} else if outFormat == "csv" {
			outputFile1 := util.GetFilenameDate("SqlServerResourceReport", "csv")
			b2, err := csvutil.Marshal(out2)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile1, string(b2))
		}

		var query3 string = `
		SELECT convert(varchar(100),name) ConfigName, convert(varchar(100),value) ConfigValue, 
		convert(varchar(100),value_in_use) ConfigValueInUse, convert(varchar(512),description) ConfigDescription
		FROM sys.configurations where configuration_id in (109,503, 505, 1532,1535,1538,1539,1543,1544,1576, 1579,1589)
		`
		fmt.Println("SQL Server Instance Configuration:")
		out3 := getInstanceConfig(connString, query3)
		if outFormat == "table" {
			olog.Print(out3)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out3))
		} else if outFormat == "csv" {
			outputFile3 := util.GetFilenameDate("SqlServerInstanceReport", "csv")
			b3, err := csvutil.Marshal(out3)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile3, string(b3))
		}

		var query4 string = `
		SELECT name as ConfigName,value as ConfigValue FROM sys.database_scoped_configurations 
		WHERE configuration_id in (1,2,4,8,13,16,18,26,35)
		`
		fmt.Println("SQL Server database scope Configuration:")
		out4 := getDbScopeConfig(connString, query4)
		if outFormat == "table" {
			olog.Print(out4)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out4))
		} else if outFormat == "csv" {
			outputFile4 := util.GetFilenameDate("SqlServerDbReport", "csv")
			b4, err := csvutil.Marshal(out4)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile4, string(b4))
		}

		var query5 string = `
		CREATE TABLE #FileSize
		(DbName NVARCHAR(128), 
			FileName NVARCHAR(128),
			PhysicalName sysname, 
			Type NVARCHAR(128),
			CurrentSizeMB DECIMAL(10,2), 
			FreeSpaceMB DECIMAL(10,2)
		);
			
		INSERT INTO #FileSize(DbName, FileName, PhysicalName, Type, CurrentSizeMB, FreeSpaceMB)
		exec sp_msforeachdb 
		'use [?]; 
		SELECT DB_NAME() AS DbName, 
			name AS FileName, 
			physical_name,
			case type_desc WHEN ''ROWS'' then ''datafile'' when ''LOG'' then ''logfile'' else type_desc end as FileType,
			size/128.0 AS CurrentSizeMB,  
			size/128.0 - CAST(FILEPROPERTY(name, ''SpaceUsed'') AS INT)/128.0 AS FreeSpaceMB
		FROM sys.database_files
		WHERE type IN (0,1);';
			
		SELECT convert(varchar(100),DbName) DbName, convert(varchar(100),FileName) FileName, convert(varchar(100),PhysicalName) PhysicalName, 
		convert(varchar(100),Type) Type, convert(varchar(100),CurrentSizeMB) CurrentSizeMB, convert(varchar(100),FreeSpaceMB) FreeSpaceMB
		FROM #FileSize
		WHERE DbName NOT IN ('distribution', 'master', 'model', 'msdb')
		`
		fmt.Println("SQL Server Database Files:")
		out5 := GetFileDetails(connString, query5)
		if outFormat == "table" {
			olog.Print(out5)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(out5))
		} else if outFormat == "csv" {
			outputFile5 := util.GetFilenameDate("SqlServerDbFileReport", "csv")
			b5, err := csvutil.Marshal(out5)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile5, string(b5))
		}

	},
}

func init() {
	sqlCmd.AddCommand(dbSqlReportCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	dbSqlReportCmd.Flags().StringP("user", "U", "", "Username to connect to SQL Server instance")
	dbSqlReportCmd.Flags().StringP("pass", "P", "", "Password to connect to SQL Server instance")
	dbSqlReportCmd.Flags().StringP("instance", "I", "", "SQL Server instance name/IP address")
	dbSqlReportCmd.Flags().Int("port", 1433, "SQL Server instance port")
	dbSqlReportCmd.Flags().String("db", "", "SQL Server database name")
	dbSqlReportCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//dbSqlReportCmd.MarkFlagRequired("server")
	//dbSqlReportCmd.MarkFlagRequired("user")
	//dbSqlReportCmd.MarkFlagRequired("pass")
}

type displySetting struct {
	ConfigName  string `json:"ConfigName"`
	ConfigValue string `json:"ConfigValue"`
}

func getDbSettingConfig(connStr, query string) []displySetting {
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
	mc1 := []displySetting{}
	for totalRows.Next() {
		c1 := displySetting{}
		err = totalRows.Scan(&c1.ConfigName, &c1.ConfigValue)
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

type instanceConfig struct {
	ConfigName        string `json:"ConfigName"`
	ConfigValue       string `json:"ConfigValue"`
	ConfigValueInUse  string `json:"ConfigValueInUse"`
	ConfigDescription string `json:"ConfigDescription"`
}
type DbScopeConfig struct {
	ConfigName  string `json:"ConfigName"`
	ConfigValue string `json:"ConfigValue"`
}

func getDbScopeConfig(connStr, query string) []DbScopeConfig {
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
	mc1 := []DbScopeConfig{}
	for totalRows.Next() {
		c1 := DbScopeConfig{}
		err = totalRows.Scan(&c1.ConfigName, &c1.ConfigValue)
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

func getInstanceConfig(connStr, query string) []instanceConfig {
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

	//row := stmt.QueryRow()
	totalRows, err := stmt.Query()
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
		panic(err)
	}
	defer totalRows.Close()

	mc1 := []instanceConfig{}
	for totalRows.Next() {
		c1 := instanceConfig{}
		err = totalRows.Scan(&c1.ConfigName, &c1.ConfigValue, &c1.ConfigValueInUse, &c1.ConfigDescription)
		if err != nil {
			panic(err)
		}

		mc1 = append(mc1, c1)
	}

	//olog.Print(mc1)

	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
	return mc1
}

type svrData struct {
	ServerSetting string `json:"srvSetting"`
	Value         string `json:"srvData"`
}

func getServerConfig(connStr, query string) []svrData {
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

	//row := stmt.QueryRow()
	totalRows, err := stmt.Query()
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
		panic(err)
	}
	defer totalRows.Close()

	mc1 := []svrData{}
	for totalRows.Next() {
		c1 := svrData{}
		err = totalRows.Scan(&c1.ServerSetting, &c1.Value)
		if err != nil {
			panic(err)
		}
		mc1 = append(mc1, c1)
	}

	// get any error encountered during iteration
	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
	//fmt.Println(util.PrettyPrint(mc1))
	//fmt.Printf("%+v", mc1)
	//olog.Print(mc1)
	return mc1
}

type FileDetails struct {
	DbName        string `json:"DbName"`
	FileName      string `json:"FileName"`
	PhysicalName  string `json:"PhysicalName"`
	Type          string `json:"Type"`
	CurrentSizeMB string `json:"CurrentSizeMB"`
	FreeSpaceMB   string `json:"FreeSpaceMB"`
}

func GetFileDetails(connStr, query string) []FileDetails {
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

	//row := stmt.QueryRow()
	totalRows, err := stmt.Query()
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
		panic(err)
	}
	defer totalRows.Close()

	var mc1 []FileDetails
	for totalRows.Next() {
		var c1 FileDetails
		err = totalRows.Scan(&c1.DbName, &c1.FileName, &c1.PhysicalName, &c1.Type, &c1.CurrentSizeMB, &c1.FreeSpaceMB)
		if err != nil {
			panic(err)
		}
		mc1 = append(mc1, c1)
	}

	// get any error encountered during iteration
	err = totalRows.Err()
	if err != nil {
		panic(err)
	}

	//olog.Print(mc1)
	return mc1
}
