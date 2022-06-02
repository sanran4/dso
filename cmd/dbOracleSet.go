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

var dbOrclSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command will apply best practices settings to oracle database layer",
	Run: func(cmd *cobra.Command, args []string) {
		srv, usr, pas, svc, prt := parseDbOrclSetFlags(cmd, args)
		connString := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", usr, pas, srv, prt, svc)
		oraAttribute, _ := cmd.Flags().GetString("attr")
		bpsFlag, _ := cmd.Flags().GetBool("bps")
		hpprocess, _ := cmd.Flags().GetBool("hpprocess")
		sgaMax, _ := cmd.Flags().GetBool("sgaMax")
		sgaTarget, _ := cmd.Flags().GetBool("sgaTarget")
		log, _ := cmd.Flags().GetBool("log")
		//RestartOracleDatabase(connString)
		sgaSize := orclGetSGAOptimal(connString)
		sgaSizeMb := sgaSize + "M"
		if oraAttribute != "" {
			attr, val := parseOracleAttribute(oraAttribute)
			//if setParameter == ""
			if attr == "_high_priority_processes" {
				setHighPriorityProcess(connString, val)
			}
			if attr == "sga_max_size" {
				setOracleSgaMax(connString, val)
			}
			if attr == "sga_target" {
				setOracleSgaTarget(connString, val)
			}
		} else {
			if bpsFlag {
				val1 := "LMS*|VKTM|LGWR"
				setHighPriorityProcess(connString, val1)
				setOracleSgaMax(connString, sgaSizeMb)
				//setOracleSgaTarget(connString, sgaSizeMb)
				val3 := "---"
				setOracleLog(connString, val3)
				fmt.Println("Applied all Oracle best practice configurations")
			} else {
				if hpprocess {
					val := "LMS*|VKTM|LGWR"
					setHighPriorityProcess(connString, val)
				}
				if sgaMax {
					setOracleSgaMax(connString, sgaSizeMb)
				}
				if sgaTarget {
					setOracleSgaTarget(connString, sgaSizeMb)
				}
				if log {
					val := "---"
					setOracleLog(connString, val)
				}
			}
		}
	},
}

func init() {
	oracleCmd.AddCommand(dbOrclSetCmd)

	dbOrclSetCmd.Flags().StringP("user", "U", "", "Username to connect to oracle instance")
	dbOrclSetCmd.Flags().StringP("pass", "P", "", "Password to connect to oracle instance")
	dbOrclSetCmd.Flags().StringP("instance", "I", "", "oracle db server IP address")
	dbOrclSetCmd.Flags().Int("port", 1521, "oracle db litioning port")
	dbOrclSetCmd.Flags().String("svc", "", "oracle db service name")
	dbOrclSetCmd.Flags().StringP("attr", "A", "", "Set individual attributes for oracle database")
	dbOrclSetCmd.Flags().Bool("bps", false, "Set all best practices for Oracle database instance at once")
	dbOrclSetCmd.Flags().Bool("hpprocess", false, "Set High Priority Process for Oracle database instance ")
	dbOrclSetCmd.Flags().Bool("sgaMax", false, "Set sga_max_size for Oracle database instance")
	dbOrclSetCmd.Flags().Bool("sgaTarget", false, "Set sga_target for Oracle database instance")
	dbOrclSetCmd.Flags().Bool("log", false, "Set best practice for log file in Oracle database")

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

	return server, user, pass, oraSvc, port
}

func setOracleSgaMax(connStr, value string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()
	//query1 := "alter system SET \"_high_priority_processes\" = '" + value + "' scope=spfile"
	query1 := "Alter system SET sga_max_size = " + value + "  scope=spfile"
	stmt1 := go_ora.NewStmt(query1, conn)
	defer stmt1.Close()
	_, err = stmt1.Exec(nil)
	dieOnError("Can't query", err)
	fmt.Println("sga_max_size has been configured.")
	fmt.Println("Restart Oracle database manually for changes to take effect...")
}

func setOracleSgaTarget(connStr, value string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()
	query2 := "alter system SET sga_target = " + value + "  scope=spfile"
	stmt2 := go_ora.NewStmt(query2, conn)
	defer stmt2.Close()
	_, err = stmt2.Exec(nil)
	dieOnError("Can't query", err)
	fmt.Println("sga_target has been configured.")
	fmt.Println("Restart Oracle database manually for changes to take effect...")
}

func setOracleLog(connStr, value string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()
	query1 := `
DECLARE 
	l_exst number(1) :=0;
	grpCount number(1) :=0;
	locGrpCnt number(3) :=0;
	testCount number(1) :=0;
	lv_statement varchar2(32676);
	currentMaxGroup number(3) :=0;
	newMaxGroup number(3) :=0;
	v_newMaxGroup varchar2(50);
	delGroup number(1) :=0;
	v_delGroup varchar2(50);
	redoLoc varchar2(50);
BEGIN
	SELECT SUBSTR(MEMBER,1,INSTR(MEMBER,'/')-1) INTO redoLoc FROM v$logfile WHERE GROUP# = (SELECT min(GROUP#) FROM v$logfile);
	SELECT max(GROUP#) INTO currentMaxGroup from v$log;
	newMaxGroup := currentMaxGroup +1;
	FOR Lcntr IN 1..5   -- LOOP FOR 5 times TO ADD 5 FILE groups
	LOOP
		SELECT count(1) INTO testCount from v$log where group# = newMaxGroup;
		IF (testCount = 0 ) THEN 
			select CAST( newMaxGroup AS varchar2(30) ) INTO v_newMaxGroup from dual;
			lv_statement := 'ALTER DATABASE ADD LOGFILE GROUP ' || v_newMaxGroup ||  '('''|| redoLoc ||''') SIZE 8G';
			execute immediate lv_statement;
		END IF;
		newMaxGroup := newMaxGroup +1;
	 END LOOP;
	dbms_lock.sleep (2); -- wait FOR 2 sec
	FOR Lcntr IN 1..120   -- LOOP FOR 10 min FOR log switching TO occur
		LOOP
			SELECT count(1) INTO grpCount from v$log where group# <= currentMaxGroup;
		  	EXIT WHEN grpCount = 0;
		  	FOR Kcntr IN 1..5
		  	LOOP 
			  	SELECT count(1) INTO locGrpCnt from v$log where group# <= currentMaxGroup and (status='INACTIVE' OR status = 'UNUSED');
			  	EXIT WHEN locGrpCnt = 0;
		  		select GROUP# into delGroup from v$log where group# <= currentMaxGroup and (status='INACTIVE' OR status = 'UNUSED') fetch first 1 row only;
		  		select CAST( delGroup AS varchar2(30) ) INTO v_delGroup from dual;
		  		lv_statement := 'ALTER DATABASE DROP LOGFILE GROUP '|| v_delGroup ;
				execute immediate lv_statement;
			  	dbms_lock.sleep (1); -- wait 1 SECOND
		  	END LOOP;
		    SELECT case when exists( select 1 from v$log where group# <= currentMaxGroup and (status='INACTIVE' OR status = 'UNUSED')) then 1 
			when exists( select 1 from v$log where group# <= currentMaxGroup and status='ACTIVE') then 2 
			when exists( select 1 from v$log where group# <= currentMaxGroup and status='CURRENT') then 3 
			else 0 end into l_exst from dual;
		  	IF (l_exst = 1 ) THEN 
		  		select GROUP# into delGroup from v$log where group# <= currentMaxGroup and (status='INACTIVE' OR status = 'UNUSED') fetch first 1 row only;
		  		select CAST( delGroup AS varchar2(30) ) INTO v_delGroup from dual;
		  		lv_statement := 'ALTER DATABASE DROP LOGFILE GROUP '|| v_delGroup ;
				execute immediate lv_statement;
			  ELSIF (l_exst = 2 ) THEN
				DBMS_OUTPUT.put_line('Waiting for log switching to occur');
			  ELSIF (l_exst = 3 ) THEN
			  	lv_statement := 'alter system switch logfile';
			  	execute immediate lv_statement;
			  ELSE
			  	DBMS_OUTPUT.put_line('Waiting for log switching to occur');
			END IF; 
			dbms_lock.sleep (4); -- wait IN SECOND
		END LOOP;
END;
	`
	//fmt.Println(query1)
	stmt2 := go_ora.NewStmt(query1, conn)
	defer stmt2.Close()
	_, err = stmt2.Exec(nil)
	dieOnError("Can't query", err)
	fmt.Println("redo log files has been modified as per the best practice. ")

}

func setHighPriorityProcess(connStr, value string) {
	conn, err := go_ora.NewConnection(connStr)
	dieOnError("Can't open the driver:", err)
	err = conn.Open()
	dieOnError("Can't open the connection:", err)
	defer conn.Close()
	query1 := "alter system SET \"_high_priority_processes\" = '" + value + "' scope=spfile"
	//fmt.Println(query1)
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
