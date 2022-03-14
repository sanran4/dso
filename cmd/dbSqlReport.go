package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var dbSqlReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Fetch report from SQL Server",
	Long:  `This report command will pull report from SQL Server`,
	Example: `dso db sql report -i 10.0.0.1 -u user1 -p pass1
dso db sql report --instance=10.0.0.1 --user=user1 --pass=pass1`,
	Run: func(cmd *cobra.Command, args []string) {

		server, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetInt("port")
		//database, _ := cmd.Flags().GetInt("database")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")

		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, pass, port)

		var query3 string = `
		SELECT convert(varchar(100),name) ConfigName, convert(varchar(100),value) ConfigValue, 
		convert(varchar(100),value_in_use) ConfigValueInUse, convert(varchar(512),description) ConfigDescription
		FROM sys.configurations where configuration_id in (109,503, 505, 1532,1535,1538,1539,1543,1544,1576, 1579,1589)
		`
		getInstanceConfig(connString, query3)

		var query4 string = `
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
			SERVERPROPERTY('IsPolyBaseInstalled') AS IsPolyBaseInstalled;
		`
		GetSrvProperty(connString, query4)

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
		GetFileDetails(connString, query5)
	},
}

func init() {
	sqlCmd.AddCommand(dbSqlReportCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	dbSqlReportCmd.Flags().StringP("user", "U", "", "Username to connect to SQL Server instance")
	dbSqlReportCmd.Flags().StringP("pass", "P", "", "Password to connect to SQL Server instance")
	dbSqlReportCmd.Flags().StringP("server", "S", "", "SQL Server instance name/IP address")
	dbSqlReportCmd.Flags().StringP("port", "p", "", "SQL Server instance port (default 1433)")
	dbSqlReportCmd.Flags().StringP("database", "d", "", "Password to connect to SQL Server instance")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	dbSqlReportCmd.MarkFlagRequired("server")
	dbSqlReportCmd.MarkFlagRequired("user")
	dbSqlReportCmd.MarkFlagRequired("pass")
}

type instanceConfig struct {
	ConfigName        string `json:"ConfigName"`
	ConfigValue       string `json:"ConfigValue"`
	ConfigValueInUse  string `json:"ConfigValueInUse"`
	ConfigDescription string `json:"ConfigDescription"`
}

func getInstanceConfig(connStr, query string) {
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

	//b, err := json.Marshal(mc1)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(string(b))
	//fmt.Println(mc1)

	fmt.Println(util.PrettyPrint(mc1))

	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
}

type svrProperty struct {
	ComputerName              string `json:"ComputerName"`
	InstanceName              string `json:"InstanceName"`
	ProductVersion            string `json:"ProductVersion"`
	ProductLevel              string `json:"ProductLevel"`
	Edition                   string `json:"Edition"`
	InstanceDefaultDataPath   string `json:"InstanceDefaultDataPath"`
	InstanceDefaultLogPath    string `json:"InstanceDefaultLogPath"`
	InstanceDefaultBackupPath string `json:"InstanceDefaultBackupPath"`
	Collation                 string `json:"Collation"`
	IsClustered               string `json:"IsClustered"`
	IsHadrEnabled             string `json:"IsHadrEnabled"`
	IsPolyBaseInstalled       string `json:"IsPolyBaseInstalled"`
}

func GetSrvProperty(connStr, query string) {
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

	mc1 := []svrProperty{}
	for totalRows.Next() {
		c1 := svrProperty{}
		err = totalRows.Scan(&c1.ComputerName, &c1.InstanceName, &c1.ProductVersion, &c1.ProductLevel, &c1.Edition, &c1.InstanceDefaultDataPath, &c1.InstanceDefaultLogPath, &c1.InstanceDefaultBackupPath, &c1.Collation, &c1.IsClustered, &c1.IsHadrEnabled, &c1.IsPolyBaseInstalled)
		if err != nil {
			panic(err)
		}
		mc1 = append(mc1, c1)
	}

	fmt.Println(util.PrettyPrint(mc1))

	//fmt.Printf("%+v", mc1)
	// get any error encountered during iteration
	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
}

type fileDetails struct {
	DbName        string `json:"DbName"`
	FileName      string `json:"FileName"`
	PhysicalName  string `json:"PhysicalName"`
	Type          string `json:"Type"`
	CurrentSizeMB string `json:"CurrentSizeMB"`
	FreeSpaceMB   string `json:"FreeSpaceMB"`
}

func GetFileDetails(connStr, query string) {
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

	mc1 := []fileDetails{}
	for totalRows.Next() {
		c1 := fileDetails{}
		err = totalRows.Scan(&c1.DbName, &c1.FileName, &c1.PhysicalName, &c1.Type, &c1.CurrentSizeMB, &c1.FreeSpaceMB)
		if err != nil {
			panic(err)
		}
		mc1 = append(mc1, c1)
	}

	fmt.Println(util.PrettyPrint(mc1))

	//fmt.Printf("%+v", mc1)
	// get any error encountered during iteration
	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
}
