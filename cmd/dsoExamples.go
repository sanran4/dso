package cmd

import (
	"fmt"
)

func db_oracle_set_ex() {
	exmpl := `
=========================== dso db oracle set examples =========================== 

EX1: // apply _high_priority_processes as per best practice recommendation
	dso db oracle set --hpprocess -I 10.0.0.1 --svc orcl -U user1 

EX2: // apply sga_max_size as per best practice recommendation
	dso db oracle set --sgaMax -I 10.0.0.2 --port 2341 --svc orcl2 -U user1 

EX3: // apply sga_target as per best practice recommendation
	dso db oracle set --sgaTarget -I 10.0.0.1 --svc orcl -U user1 

EX4: // apply all best practices as per the recommendation
	dso db oracle set --bps -I 10.0.0.1 --svc orcl -U user1 

EX5 //  apply individual best practice with custom value
	dso db oracle set -A "sga_max_size=24G" -I 10.0.0.1 --svc orcl -U user1 -P pass1
	or
	dso db oracle set -A "_high_priority_processes=LMS*|VKTM|LGWR" -I 10.0.0.1 --svc orcl -U user1 -P pass1

EX6: // apply redo log file best practice recommendation
	dso db oracle set --log -I 10.0.0.1 --svc orcl -U user1 

===================================================================================
	`
	fmt.Println(exmpl)
}

func db_oracle_get_ex() {
	exmpl := `
=========================== dso db oracle get examples =========================== 

EX1: // validate best practice recommendation for oracle database (output in TABLE format)
	dso db oracle get  -I 10.0.0.1 --svc orcl -U user1 

EX2: // validate best practice recommendation for oracle database (output in JSON format)
	dso db oracle get  -I 10.0.0.1 --svc orcl -U user1 -o json

EX3: // validate best practice recommendation for oracle database (output in csv format)
	dso db oracle get  -I 10.0.0.1 --svc orcl -U user1 -o csv

EX4: // validate best practice recommendation for oracle database on non default port
	dso db oracle get -I 10.0.0.2 --port 2341 --svc orcl2 -U user1 

EX5: // validate best practice recommendation for oracle database with inline password
	dso db oracle get -I 10.0.0.2 --svc orcl2 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func db_oracle_report_ex() {
	exmpl := `
=========================== dso db oracle report examples =========================== 

EX1: // fetch report from oracle database (output in TABLE format)
	dso db oracle report  -I 10.0.0.1 --svc orcl -U user1 

EX2: // fetch report from oracle database (output in JSON format)
	dso db oracle report  -I 10.0.0.1 --svc orcl -U user1 -o json

EX3: // fetch report from oracle database (output in csv format)
	dso db oracle report  -I 10.0.0.1 --svc orcl -U user1 -o csv

EX4: // fetch report from oracle database on non default port
	dso db oracle report -I 10.0.0.2 --port 2341 --svc orcl2 -U user1 

EX5: // fetch report from oracle database with inline password
	dso db oracle report -I 10.0.0.2 --svc orcl2 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func db_sql_get_ex() {
	exmpl := `
=========================== dso db sql get examples =========================== 

EX1: // validate best practice recommendation from SQL Server database (output in TABLE format)
	dso db sql get  -I 10.0.0.1 -U user1 

EX2: // validate best practice recommendation from SQL Server database (output in JSON format)
	dso db sql get  -I 10.0.0.1 -U user1 -o json

EX3: // validate best practice recommendation from SQL Server database (output in CSV format)
	dso db sql get  -I 10.0.0.1 -U user1 -o csv

EX4: // validate best practice recommendation from SQL Server database on non default port
	dso db sql get -I 10.0.0.2 --port 2341 -U user1 

EX5: // validate best practice recommendation from SQL Server database with inline password
	dso db sql get -I 10.0.0.2 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func db_sql_report_ex() {
	exmpl := `
=========================== dso db sql report examples =========================== 

EX1: // fetch report from SQL Server database (output in TABLE format)
	dso db sql report  -I 10.0.0.1 -U user1 

EX2: // fetch report from SQL Server database (output in JSON format)
	dso db sql report  -I 10.0.0.1 -U user1 -o json

EX3: // fetch report from SQL Server database (output in CSV format)
	dso db sql report  -I 10.0.0.1 -U user1 -o csv

EX4: // fetch report from SQL Server database on non default port
	dso db sql report -I 10.0.0.2 --port 2341 -U user1 

EX5: // fetch report from SQL Server database with inline password
	dso db sql report -I 10.0.0.2 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func db_sql_set_ex() {
	exmpl := `
=========================== dso db sql set examples =========================== 

EX1: // apply CPU affinity mask as per best practice recommendation
	dso db sql set --affinity -I 10.0.0.1 -U user1 

EX2: // apply Min and Max memory settings as per best practice recommendation
	dso db sql set --sqlmem -I 10.0.0.1 -U user1 

EX3: // apply db file auto-growth as per best practice recommendation
	dso db sql set --autogrowth -I 10.0.0.1 -U user1 

EX4: // apply Query Optimizer HotFixes as per best practice recommendation
	dso db sql set --qohf -I 10.0.0.1 -U user1 

EX5: // apply Dedicated remote admin connection as per best practice recommendation
	dso db sql set --dac -I 10.0.0.1 -U user1 

EX6: // apply all best practices per best practice recommendation
	dso db sql set --bps -I 10.0.0.1 -U user1

EX7: // apply best practice recommendation with inline password
	dso db sql set --bps -I 10.0.0.1 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func storage_pstore_report_ex() {
	exmpl := `
=========================== dso storage pstore report examples =========================== 

EX1: // fetch report from PowerStore storage (output in TABLE format)
	dso storage pstore report -I 10.0.0.1 -U user1 

EX2: // fetch report from PowerStore storage (output in JSON format)
	dso storage pstore report  -I 10.0.0.1 -U user1 -o json

EX3: // fetch report from PowerStore storage (output in CSV format)
	dso storage pstore report -I 10.0.0.1 -U user1 -o csv

EX4: // fetch report from PowerStore storage for specific host group
	dso storage pstore report --hgroup -I 10.0.0.2 -U user1 

EX5: // fetch report from PowerStore storage with inline password
	dso storage pstore report -I 10.0.0.2 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func os_rhel_get_ex() {
	exmpl := `
=========================== dso os rhel get examples =========================== 

EX1: // validate all best practice recommendation on RedHat Linux OS for oracle workload (output in TABLE format)
	dso os rhel get -w oracle --bps -I 10.0.0.1 -U user1 

EX2: // validate kernal parameters best practice recommendation on RedHat Linux OS for oracle workload (output in JSON format)
	dso os rhel get -w oracle --tunedadm -I 10.0.0.1 -U user1 -o json

EX3: // validate Hugepages best practice recommendation on RedHat Linux OS for oracle workload (output in CSV format)
	dso os rhel get -w oracle --hpage -I 10.0.0.1 -U user1 -o csv

EX4: // validate all best practice recommendation on RedHat Linux OS for SQL Server workload (output in TABLE format)
	dso os rhel get -w sql --bps -I 10.0.0.1 -U user1 

EX5: // validate kernal parameters best practice recommendation on RedHat Linux OS for SQL Server workload (output in JSON format)
	dso os rhel get -w sql --tunedadm -I 10.0.0.1 -U user1 -o json

EX6: // validate disk related best practice recommendation on RedHat Linux OS for SQL Server workload (output in CSV format)
	dso os rhel get -w sql --disk -I 10.0.0.1 -U user1 -o csv

EX7: // validate MSSQL Configuration related best practice recommendation on RedHat Linux OS for SQL Server workload
	dso os rhel get -w sql --msconf -I 10.0.0.1 -U user1 -o csv

EX8: // validate kernal parameters and MSSQL Configuration related best practice recommendation on RedHat Linux OS for SQL Server workload 
	dso os rhel get -w sql --tunedadm --msconf -I 10.0.0.1 -U user1

EX9: // validate best practice recommendation on RedHat Linux OS for SQL Server workload with inline password
	dso os rhel get -w sql --bps -I 10.0.0.2 -U user1 -P pass1

EX10: // validate best practice recommendation on RedHat Linux OS for SQL Server workload with non default ssh port
	dso os rhel get -w sql --bps -I 10.0.0.2 --portSSH 2222 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func os_rhel_report_ex() {
	exmpl := `
=========================== dso os rhel report examples =========================== 

EX1: // fetch report from RedHat Linux OS for oracle workload (output in TABLE format)
	dso os rhel report -w oracle -I 10.0.0.1 -U user1 

EX2: // fetch report from RedHat Linux OS for SQL Server workload (output in JSON format)
	dso os rhel report -w sql -I 10.0.0.1 -U user1 -o json

EX3: // fetch report from RedHat Linux OS for oracle workload (output in CSV format)
	dso os rhel report -w oracle -I 10.0.0.1 -U user1 -o json

EX4: // fetch report from RedHat Linux OS for oracle workload with inline password
	dso os rhel report -w oracle -I 10.0.0.1 -U user1 -P pass1

EX5: // fetch report from RedHat Linux OS for SQL Server workload with non default ssh port
	dso os rhel report -w oracle -I 10.0.0.1 --portSSH 2222 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func os_rhel_set_ex() {
	exmpl := `
	=========================== dso os rhel set examples =========================== 
	
	EX01: // apply all best practice recommendation on RedHat Linux OS for oracle workload 
		dso os rhel set -w oracle --bps -I 10.0.0.1 -U user1 
	
	EX02: // apply kernal parameters best practice recommendation on RedHat Linux OS for oracle workload 
		dso os rhel set -w oracle --tunedadm -I 10.0.0.1 -U user1 
	
	EX03: // apply Hugepages best practice recommendation on RedHat Linux OS for oracle workload 
		dso os rhel set -w oracle --hpage -I 10.0.0.1 -U user1 
	
	EX04: // apply all best practice recommendation on RedHat Linux OS for SQL Server workload 
		dso os rhel set -w sql --bps -I 10.0.0.1 -U user1 
	
	EX05: // apply kernal parameters best practice recommendation on RedHat Linux OS for SQL Server workload 
		dso os rhel set -w sql --tunedadm -I 10.0.0.1 -U user1 
	
	EX06: // apply disk related best practice recommendation on RedHat Linux OS for SQL Server workload 
		dso os rhel set -w sql --disk -I 10.0.0.1 -U user1 
	
	EX07: // apply MSSQL Configuration related best practice recommendation on RedHat Linux OS for SQL Server workload
		dso os rhel set -w sql --msconf -I 10.0.0.1 -U user1 
	
	EX08: // apply kernal parameters and MSSQL Configuration related best practice recommendation on RedHat Linux OS for SQL Server workload 
		dso os rhel set -w sql --tunedadm --msconf -I 10.0.0.1 -U user1
	
	EX09: // apply best practice recommendation on RedHat Linux OS for SQL Server workload with inline password
		dso os rhel set -w sql --bps -I 10.0.0.2 -U user1 -P pass1
	
	EX10: // apply best practice recommendation on RedHat Linux OS for SQL Server workload with non default ssh port
		dso os rhel set -w sql --bps -I 10.0.0.2 --portSSH 2222 -U user1 -P pass1
	
	Ex11:- // apply custom SQL Server memory limit using mssql-conf on RedHat Linux OS for SQL Server workload
		dso os rhel set -w sql -A "memory.memorylimitmb=8192" -I 10.0.0.1 -U user1 -P pass1

	Ex12:- // enable custom SQL Server traceflag
		dso os rhel set -w sql -A "traceflag=834" -I 10.0.0.1 -U user1 -P pass1 

	===================================================================================
		`
	fmt.Println(exmpl)
}

func server_get_ex() {
	exmpl := `
=========================== dso server get examples =========================== 

EX1: // validate intel based server BIOS best practice recommendation for database workload (output in TABLE format)
	dso server get --bios -I 10.0.0.1 -U user1 

EX2: // validate intel based server BIOS best practice recommendation for database workload (output in JSON format)
	dso server get --bios -I 10.0.0.1 -U user1  -o json

EX3: // validate intel based server BIOS best practice recommendation for database workload (output in CSV format)
	dso server get --bios -I 10.0.0.1 -U user1  -o csv

EX4: // validate intel based server BIOS best practice recommendation for database workload with inline password
	dso server get --bios -I 10.0.0.1 -U user1 -P pass1

EX5: // check BIOS config job status based on job_id on intel based server  
	dso server get --jobStatus=JID_526269044866 -I 10.0.0.1 -U user1 -P pass1
	or
	dso server get -j JID_526269044866 -I 10.0.0.1 -U user1 

===================================================================================
	`
	fmt.Println(exmpl)
}

func server_report_ex() {
	exmpl := `
=========================== dso server report examples =========================== 

EX1: // fetch report from intel based server BIOS (output in TABLE format)
	dso server report -I 10.0.0.1 -U user1 

EX2: // fetch report from intel based server BIOS (output in JSON format)
	dso server report -I 10.0.0.1 -U user1  -o json

EX3: // fetch report from intel based server BIOS (output in CSV format)
	dso server report -I 10.0.0.1 -U user1  -o csv

EX4: // fetch report from intel based server BIOS with inline password
	dso server report -I 10.0.0.1 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}

func server_set_ex() {
	exmpl := `
=========================== dso db sql set examples =========================== 

EX01: // Apply System Profile on intel based server BIOS 
	dso server set -A "SysProfile:PerfOptimized" -I 10.0.0.1 -U user1 

EX02: // Enable Processor Virtualization on intel based server BIOS 
	dso server set -A "ProcVirtualization:Enabled" -I 10.0.0.1 -U user1 -P pass1 

EX03: // Enable ProcX2Apic for processor on intel based server BIOS 
	dso server set -A "ProcX2Apic:Enabled" -I 10.0.0.1 -U user1 

EX04: // Enable Logical Processor/threads on intel based server BIOS 
	dso server set -A "LogicalProc:Enabled" -I 10.0.0.1 -U user1 -P pass1

EX05: // Apply Memory Operating Mode  on intel based server BIOS 
	dso server set -A "MemOpMode:OptimizerMode" -I 10.0.0.1 -U user1 

EX06: // Disable Serial Communication on intel based server BIOS 
	dso server set -A "SerialComm:Off" -I 10.0.0.1 -U user1 -P pass1

EX07: // Disable USB Ports on intel based server BIOS 
	dso server set -A "ProcVirtualization:Enabled" -I 10.0.0.1 -U user1 -P pass1

EX08: // Disable USB Managed Port attribute on intel based server BIOS 
	dso server set -A "UsbManagedPort:Off" -I 10.0.0.1 -U user1 

EX09: // apply Workload Profile on intel based server BIOS 
	dso server set -A "WorkloadProfile:DbOptimizedProfile" -I 10.0.0.1 -U user1 -P pass1

EX10: // apply all best practice reccomendation on intel based server BIOS 
	dso server set --bps -I 10.0.0.1 -U user1 -P pass1

EX11: // create BIOS config job on intel based server  
	dso server set --bps -I 10.0.0.1 -U user1 -P pass1

EX12: // reboot intel based server  
	dso server set --bps -I 10.0.0.1 -U user1 -P pass1

EX13: // check BIOS config job status (in continuous loop) based on job_id on intel based server  
	dso server set --jobStatus=JID_526269044866 -I 10.0.0.1 -U user1 -P pass1
	or
	dso server set -j JID_526269044866 -I 10.0.0.1 -U user1 -P pass1

===================================================================================
	`
	fmt.Println(exmpl)
}
