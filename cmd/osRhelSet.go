package cmd

import (
	"fmt"
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
	Short: "Set best practice settings values",
	Long:  `This set command will set Dell Recomended best practice values for SQL/Oracle  solution on RHEL OS`,
	Example: `
Ex1:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -t // set Tuned-Adm profile for SQL Server workload
Ex2:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 --tunedAdm // set Tuned-Adm profile for SQL Server workload
Ex3:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -w sql -t // set Tuned-Adm profile for SQL Server workload
Ex4:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -w sql -m // set MSSQL-CONF best practice for SQL Server workload
Ex5:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 --msConf // set MSSQL-CONF best practice for SQL Server workload
Ex6:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -w sql -t -m // set Tuned-Adm & MSSQL-CONF best practice for SQL Server workload
Ex7:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -w sql -A "memory.memorylimitmb=8192" // set SQL Server memory limit using mssql-conf
Ex8:- dso os rhel set -I 10.0.0.1 -U user1 -P pass1 -w sql -A "traceflag=834" // set SQL Server traceflag
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
						restartSQL(c)
					}
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
	osRhelSetCmd.Flags().StringP("workload", "w", "sql", "Application workload (sql/oracle)")
	osRhelSetCmd.Flags().BoolP("tunedAdm", "t", false, "Set settings for optimal tuned-Adm profile")
	osRhelSetCmd.Flags().BoolP("msConf", "m", false, "Set setting for optimal mssql-conf")
	osRhelSetCmd.Flags().StringP("attr", "A", "", "Set individual attributes (ex:- -A \"memory.memorylimitmb=8192\") please help for more info.")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	osRhelSetCmd.MarkFlagRequired("ip")
	osRhelSetCmd.MarkFlagRequired("user")
	osRhelSetCmd.MarkFlagRequired("pass")
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
	cmd1 := `systemctl restart mssql-server`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
}

func initOsRhelSetStep(cmd *cobra.Command, args []string) (*ssh.Client, error) {

	ip, _ := cmd.Flags().GetString("ip")
	portSSH, _ := cmd.Flags().GetString("portSSH")
	user, _ := cmd.Flags().GetString("user")
	pass, _ := cmd.Flags().GetString("pass")
	workload, _ = cmd.Flags().GetString("workload")
	attribute, _ = cmd.Flags().GetString("attr")

	tunedAdm, _ = cmd.Flags().GetBool("tunedAdm")
	mssqlConf, _ = cmd.Flags().GetBool("msConf")

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
		panic("Failed to dial: " + err.Error())
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
	time.Sleep(2 * time.Second)

	cmd2 := `
systemctl restart mssql-server.service
`
	_, err = util.ExecCmd(client, cmd2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Change completed successfully... ")
}
