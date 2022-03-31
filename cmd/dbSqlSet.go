package cmd

import (
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var dbSqlSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command will configure Dell Recomended bast practice settings to SQL Server",
	Long:  `This set command will configure Dell Recomended bast practice settings to SQL Server`,
	Example: `
EX1: dso db sql set -S 10.0.0.1 -U user1 
EX2: dso db sql set -S 10.0.0.1 -U user1 -P pass1
EX3: dso db sql set --server=10.0.0.1 --user=user1 --pass=pass1
`,
	Run: func(cmd *cobra.Command, args []string) {

		server, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetInt("port")
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

		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, pass, port)

		var query2 string = `
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
		getServerConfig(connString, query2)

	},
}

func init() {
	sqlCmd.AddCommand(dbSqlSetCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	dbSqlSetCmd.Flags().StringP("user", "U", "", "Username to connect to SQL Server instance")
	dbSqlSetCmd.Flags().StringP("pass", "P", "", "Password to connect to SQL Server instance")
	dbSqlSetCmd.Flags().StringP("server", "S", "", "SQL Server instance name/IP address")
	dbSqlSetCmd.Flags().Int("port", 1433, "SQL Server instance port")
	dbSqlSetCmd.Flags().String("db", "", "SQL Server database name")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	dbSqlSetCmd.MarkFlagRequired("server")
	dbSqlSetCmd.MarkFlagRequired("user")
	//dbSqlReportCmd.MarkFlagRequired("pass")
}

/*
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

	//fmt.Println(util.PrettyPrint(mc1))
	olog.Print(mc1)

	err = totalRows.Err()
	if err != nil {
		panic(err)
	}
}

type svrData struct {
	ServerSetting string `json:"srvSetting"`
	Value         string `json:"srvData"`
}

func getServerConfig(connStr, query string) {
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
	olog.Print(mc1)
}

type FileDetails struct {
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

	//fmt.Println(mc1)
	//fmt.Println(util.PrettyPrint(mc1))
	//tableprinter.Print(os.Stdout, mc1)
	//b, err := json.Marshal(mc1)
	//if err != nil {
	//	panic(err)
	//}
	//tableprinter.PrintJSON(os.Stdout, b)
	olog.Print(mc1)
}
*/
