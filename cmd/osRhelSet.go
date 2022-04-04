package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

//var tunedAdm bool
//var mssqlConf bool
var workload string = "sql"
var attribute string = ""

var osRhelSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This set command with configure Dell Recomended best practice settings at RHEL operating system layer",
	Long:  `This set command for RHEL sub-module will set Dell Recomended best practice values for SQL/Oracle solution on RHEL OS`,
	Example: `
Ex1:- dso os rhel set --tunedadm -I 10.0.0.1 -U user1 -P pass1                // set Tuned-Adm profile for SQL Server workload
Ex2:- dso os rhel set -w sql --tunedadm -I 10.0.0.1 -U user1 -P pass1         // set Tuned-Adm profile for SQL Server workload
Ex3:- dso os rhel set -w sql --msconf -I 10.0.0.1 -U user1 -P pass1           // set MSSQL-CONF best practice for SQL Server workload
Ex4:- dso os rhel set --msconf -I 10.0.0.1 -U user1 -P pass1                  // set MSSQL-CONF best practice for SQL Server workload
Ex5:- dso os rhel set --tunedadm --msconf-I 10.0.0.1 -U user1 -P pass1        // set both Tuned-Adm & MSSQL-CONF best practices for SQL Server workload
Ex6:- dso os rhel set -w sql -A "memory.memorylimitmb=8192" -I 10.0.0.1 -U user1 -P pass1  // set SQL Server memory limit using mssql-conf
Ex7:- dso os rhel set -A "memory.memorylimitmb=8192" -I 10.0.0.1 -U user1 -P pass1  // set SQL Server memory limit using mssql-conf
Ex8:- dso os rhel set -w sql -A "traceflag=834" -I 10.0.0.1 -U user1 -P pass1 // set SQL Server traceflag
Ex9:- dso os rhel set -A "traceflag=834" -I 10.0.0.1 -U user1 -P pass1        // set SQL Server traceflag
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := initOsRhelSetStep(cmd, args)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		if workload == "sql" {
			var usrCnf string = "n"
			if attribute != "" {
				att, val := parseAttribute(attribute)
				fmt.Println("This action require SQL Server instance restart. Are you sure want to continue? (y/n): ")
				fmt.Scanln(&usrCnf)
				if strings.ToLower(usrCnf) == "y" {
					setAttrSetting(c, att, val)
					restartSQL(c)
				}
			} else {
				if tunedAdm {
					setTunedAdmSettings(c)
				}
				if mssqlConf {
					fmt.Println("This action require SQL Server instance restart. Are you sure want to continue? (y/n): ")
					fmt.Scanln(&usrCnf)
					if strings.ToLower(usrCnf) == "y" {
						setMsConfigSettings(c)
						fmt.Println("Restarting SQL Server instance... ")
						restartSQL(c)
					}
				}
				if !tunedAdm && !mssqlConf {
					fmt.Println("no sub flag (--tunedadm or --msconf) provided")
					fmt.Println("use below instruction to see help and examples for this command")
					fmt.Println("dso os rhel set --help")
				}
			}
		}
	},
}

func init() {
	rhelCmd.AddCommand(osRhelSetCmd)
	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	osRhelSetCmd.Flags().StringP("ip", "I", "", "IP / FQDN of the RHEL operating system")
	osRhelSetCmd.Flags().StringP("portSSH", "p", "22", "SSH port for connecting to RHEL os")
	osRhelSetCmd.Flags().StringP("user", "U", "", "Username for the RHEL operating system")
	osRhelSetCmd.Flags().StringP("pass", "P", "", "Password for the RHEL operating system")
	osRhelSetCmd.Flags().StringP("workload", "w", "sql", "Application workload [sql/oracle]")
	osRhelSetCmd.Flags().Bool("tunedadm", false, "Set settings for optimal tuned-Adm profile")
	osRhelSetCmd.Flags().Bool("msconf", false, "Set setting for optimal mssql-conf")
	osRhelSetCmd.Flags().StringP("attr", "A", "", "Set individual attributes for mssql-conf(ex: -A \"memory.memorylimitmb=8192\") see help for more info.")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//osRhelSetCmd.MarkFlagRequired("ip")
	//osRhelSetCmd.MarkFlagRequired("user")
	//osRhelSetCmd.MarkFlagRequired("pass")
}

func parseAttribute(str string) (attr, val string) {
	tmp := strings.Split(str, "=")
	attr = tmp[0]
	val = tmp[1]
	return
}

func setAttrSetting(client *ssh.Client, attr, val string) error {
	var cmnd string = ""
	cmnd = "/opt/mssql/bin/mssql-conf"
	if attr == "traceflag" {
		cmnd = cmnd + " traceflag " + val + " on"
	} else {
		cmnd = cmnd + " set " + attr + " " + val
	}

	fmt.Println(cmnd)
	_, err := util.ExecCmd(client, cmnd)
	if err != nil {
		panic(err)
	}
	return nil
}

func restartSQL(client *ssh.Client) {
	time.Sleep(1 * time.Second)
	cmd1 := `systemctl restart mssql-server.service`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	fmt.Println("SQL Server restarted successfully... ")
}

func initOsRhelSetStep(cmd *cobra.Command, args []string) (*ssh.Client, error) {

	ip, ok := os.LookupEnv("RHEL_OS_HOST")
	if !ok {
		ip, _ = cmd.Flags().GetString("ip")
	}
	user, ok := os.LookupEnv("RHEL_OS_USER")
	if !ok {
		user, _ = cmd.Flags().GetString("user")
	}
	//ip, _ := cmd.Flags().GetString("ip")
	portSSH, _ := cmd.Flags().GetString("portSSH")
	//user, _ := cmd.Flags().GetString("user")
	pass, _ := cmd.Flags().GetString("pass")
	var err error
	if pass == "" {
		pass, err = util.GetPasswd()
		if err != nil {
			log.Printf("error getting password %v", err)
		}
	}
	workload, _ = cmd.Flags().GetString("workload")
	attribute, _ = cmd.Flags().GetString("attr")

	tunedAdm, _ = cmd.Flags().GetBool("tunedadm")
	mssqlConf, _ = cmd.Flags().GetBool("msconf")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ipAddr := ip + ":" + portSSH
	client, err := ssh.Dial("tcp", ipAddr, config)
	if err != nil {
		//panic("Failed to dial: " + err.Error())
		return nil, err
	}
	//defer client.Close()

	return client, nil
}

func setTunedAdmSettings(client *ssh.Client) {
	cmd1 := `
	mkdir -p /usr/lib/tuned/mssql
	`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())

	cmd2 := `
cat > /usr/lib/tuned/mssql/tuned.conf << "EOF"
#
# A Tuned configuration for SQL Server on Linux
#

[main]
summary=Optimize for Microsoft SQL Server
include=throughput-performance

[cpu]
force_latency=5

[sysctl]
vm.swappiness = 1
vm.dirty_background_ratio = 3
vm.dirty_ratio = 80
vm.dirty_expire_centisecs = 500
vm.dirty_writeback_centisecs = 100
vm.transparent_hugepages=always
# For multi-instance SQL deployments, use
# vm.transparent_hugepages=madvise
vm.max_map_count=1600000
net.core.rmem_default = 262144
net.core.rmem_max = 4194304
net.core.wmem_default = 262144
net.core.wmem_max = 1048576
kernel.numa_balancing=0
kernel.sched_min_granularity_ns = 15000000
kernel.sched_wakeup_granularity_ns = 2000000
EOF
`
	_, err = execCmd(client, cmd2)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res2.String())

	cmd3 := `
chmod +x /usr/lib/tuned/mssql/tuned.conf
tuned-adm profile mssql
tuned-adm active
`
	res3, err := execCmd(client, cmd3)
	if err != nil {
		panic(err)
	}
	fmt.Println(res3.String())
}

func setMsConfigSettings(client *ssh.Client) {
	cmd1 := `
	/opt/mssql/bin/mssql-conf set control.alternatewritethrough 0 && /opt/mssql/bin/mssql-conf set control.writethrough 1 && /opt/mssql/bin/mssql-conf traceflag 3979 834 on
	`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	//time.Sleep(2 * time.Second)

	//cmd2 := `systemctl restart mssql-server.service`
	//_, err = util.ExecCmd(client, cmd2)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("mssql-conf changes applied successfully... ")
}
