package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var dbSqlSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command will apply bast practice settings to SQL Server",
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
		bpsFlag, _ := cmd.Flags().GetBool("bps")
		affinity, _ := cmd.Flags().GetBool("affinity")
		sqlmem, _ := cmd.Flags().GetBool("sqlmem")
		autogrowth, _ := cmd.Flags().GetBool("autogrowth")
		qohf, _ := cmd.Flags().GetBool("qohf")
		dac, _ := cmd.Flags().GetBool("dac")

		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", server, user, pass, port)

		if affinity || bpsFlag {
			sqlSetAffinityMask(connString)
			fmt.Println("applied Affinity mask as per best practice..")
		}
		if sqlmem || bpsFlag {
			sqlSetSqlMemory(connString)
			fmt.Println("applied Min and Max memory settings as per best practice..")
		}
		if autogrowth || bpsFlag {
			sqlSetFileAutoGrowth(connString)
			fmt.Println("applied File auto-growth as per best practice..")
		}
		if qohf || bpsFlag {
			sqlSetQueryOptimizerHotFixes(connString)
			fmt.Println("applied Query Optimizer HotFixes as per best practice..")
		}
		if dac || bpsFlag {
			sqlsetDAC(connString)
			fmt.Println("applied Dedicated remote admin connection as per best practice..")
		}

	},
}

func init() {
	sqlCmd.AddCommand(dbSqlSetCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	dbSqlSetCmd.Flags().StringP("user", "U", "", "Username to connect to SQL Server instance")
	dbSqlSetCmd.Flags().StringP("pass", "P", "", "Password to connect to SQL Server instance")
	dbSqlSetCmd.Flags().StringP("instance", "I", "", "SQL Server instance name/IP address")
	dbSqlSetCmd.Flags().Int("port", 1433, "SQL Server instance port")
	dbSqlSetCmd.Flags().String("db", "", "SQL Server database name")
	dbSqlSetCmd.Flags().Bool("bps", false, "Set SQL Server all best practices at once")
	dbSqlSetCmd.Flags().Bool("affinity", false, "Set SQL Server CPU Affinity at instance level")
	dbSqlSetCmd.Flags().Bool("sqlmem", false, "Set SQL Server Min & Max Memory at instance level")
	dbSqlSetCmd.Flags().Bool("autogrowth", false, "Set SQL Server data and log file growth for user databases")
	dbSqlSetCmd.Flags().Bool("qohf", false, "Set SQL Server Query Optimizer HotFixes for user databases")
	dbSqlSetCmd.Flags().Bool("dac", false, "Enable SQL Server Dedicated remote admin connection")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//dbSqlSetCmd.MarkFlagRequired("server")
	//dbSqlSetCmd.MarkFlagRequired("user")
	//dbSqlReportCmd.MarkFlagRequired("pass")
}

func sqlSetAffinityMask(connStr string) {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	var query string = `
declare @cmd nvarchar(1000)
select @cmd = 'ALTER SERVER CONFIGURATION SET PROCESS AFFINITY CPU = 0 TO '+( convert(varchar(10),(cpu_count-1)) ) from sys.dm_os_sys_info
EXECUTE sp_executesql @cmd
	`
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("SQL Server CPU Affinity set successfully")
}

func sqlSetSqlMemory(connStr string) {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	var query string = `
	declare @currSqlMem int, @sqlMem int
	select @currSqlMem = (total_physical_memory_kb/1024) from sys.dm_os_sys_memory
	select @currSqlMem
	if exists (select 1 FROM sys.dm_os_host_info where host_platform = 'Linux')
		set @sqlMem = @currSqlMem
	ELSE
		set @sqlMem = ((@currSqlMem * 0.8)/2)*2
	;
	BEGIN
	exec sp_configure 'show advanced options', 1;
	reconfigure
	exec sp_configure 'min server memory', @sqlMem;
	exec sp_configure 'max server memory', @sqlMem;
	exec sp_configure 'show advanced options', 0;
	reconfigure
	END
	`
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("SQL Server Min and Max memory set successfully")
}

func sqlSetFileAutoGrowth(connStr string) {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	var query string = `
declare @growthSizeInMb int, @cmd nvarchar(max)
set @growthSizeInMb = 1024
select @cmd = ''
select @cmd = @cmd +'; '+'ALTER DATABASE '+db_name(database_id)+'
	MODIFY FILE ( NAME = N'''+name+''', FILEGROWTH = '+ convert(nvarchar(100),@growthSizeInMb)+'MB )'
from sys.master_files 
where database_id not in (db_id('master'),db_id('model'),db_id('msdb'),db_id('tempdb'))
select @cmd = right(@cmd, len(@cmd)-1)
EXECUTE sp_executesql @cmd
	`
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("SQL Server Database filegrowth set successfully")
}

func sqlSetQueryOptimizerHotFixes(connStr string) {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	var query string = `
	DECLARE @dbName TABLE (dbname sysname)
	declare @cnt int
	insert into @dbName(dbname)
	select name from sys.sysdatabases where name not in ('master','model','msdb','tempdb')
	select @cnt = count(1) from @dbName
	while @cnt > 0
	BEGIN
	declare @cmd nvarchar(1000)
	declare @dname nvarchar(255)
	select top 1 @dname = dbname from @dbName
	select top 1 @cmd = 'use '+ @dname + '; ALTER DATABASE SCOPED CONFIGURATION SET QUERY_OPTIMIZER_HOTFIXES = ON;' 
	EXECUTE sp_executesql @cmd
	delete from @dbName where dbname = @dname
	select @cnt = count(1) from @dbName
	END
	`
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("SQL Server QUERY_OPTIMIZER_HOTFIXES set successfully")
}

func sqlsetDAC(connStr string) {
	conn, err := sql.Open("mssql", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	var query string = `
	BEGIN
	exec sp_configure 'show advanced options', 1;
	exec sp_configure 'remote admin connections', 1;
	exec sp_configure 'show advanced options', 0;
	reconfigure
	END
	`
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("SQL Server Dedicated Remote Admin Connection(DAC) has been enabled successfully")
}
