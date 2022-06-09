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
	Short: "This set command with apply best practice settings at RHEL operating system",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := initOsRhelSetStep(cmd, args)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		force, _ := cmd.Flags().GetBool("force")
		disk, _ := cmd.Flags().GetBool("disk")
		bps, _ := cmd.Flags().GetBool("bps")
		var usrCnf string = "n"
		if workload == "sql" {
			if attribute != "" {
				att, val := parseAttribute(attribute)
				if !force {
					fmt.Println("This action require SQL Server instance restart. Are you sure want to continue? (y/n): ")
					fmt.Scanln(&usrCnf)
					if strings.ToLower(usrCnf) == "y" {
						setAttrSetting(c, att, val)
						restartSQL(c)
					}
				} else {
					setAttrSetting(c, att, val)
					fmt.Println("Restarting SQL Server instance... ")
					restartSQL(c)
				}
			} else {
				if tunedAdm || bps {
					if !force {
						fmt.Println("This action will change kernel parameters. Are you sure want to continue? (y/n): ")
						fmt.Scanln(&usrCnf)
						if strings.ToLower(usrCnf) == "y" {
							setTunedAdmSettingsSql(c)
						}
					} else {
						setTunedAdmSettingsSql(c)
					}
				}
				if disk || bps {
					fmt.Println("Setting disk related best practices for SQL Server workload ")
					setMssqlDiskSettings(c)
				}
				if mssqlConf || bps {
					if !force {
						fmt.Println("This action require SQL Server instance restart. Are you sure want to continue? (y/n): ")
						fmt.Scanln(&usrCnf)
						if strings.ToLower(usrCnf) == "y" {
							setMsConfigSettings(c)
							fmt.Println("Restarting SQL Server instance... ")
							restartSQL(c)
						}
					} else {
						setMsConfigSettings(c)
						fmt.Println("Restarting SQL Server instance... ")
						restartSQL(c)
					}
				}
				if !tunedAdm && !disk && !mssqlConf && !bps {
					fmt.Println("no sub flag (--bps or --tunedadm or --disk or --msconf) provided")
					fmt.Println("use below instruction to see help and examples for this command")
					fmt.Println("dso os rhel -w sql set --help")
				}
			}
		}
		if workload == "oracle" {
			if attribute != "" {
				att, val := parseAttribute(attribute)
				if att == "hugepages" {
					var cmd string
					cmd = "sysctl -w vm.nr_hugepages="
					if val == "" {
						res := getHugepagesRecomendValue(c)
						//fmt.Println(res)
						cmd = cmd + res
						setHugePages(c, cmd)
					} else {
						cmd = cmd + val
						setHugePages(c, cmd)
					}
				}
			} else {
				if tunedAdm {
					if !force {
						fmt.Println("This action will change kernel parameters. Are you sure want to continue? (y/n): ")
						fmt.Scanln(&usrCnf)
						if strings.ToLower(usrCnf) == "y" {
							setTunedAdmSettingsOrcl(c)
						}
					} else {
						setTunedAdmSettingsOrcl(c)
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
	osRhelSetCmd.Flags().StringP("workload", "w", "", "Application workload [sql/oracle]")
	osRhelSetCmd.Flags().Bool("tunedadm", false, "Set settings for optimal tuned-Adm profile")
	osRhelSetCmd.Flags().Bool("msconf", false, "Set setting for optimal mssql-conf")
	osRhelSetCmd.Flags().Bool("disk", false, "Set disk related best practice for sql server workload on OS layer")
	osRhelSetCmd.Flags().Bool("bps", false, "Set all best practices for given workload on OS layer")
	osRhelSetCmd.Flags().Bool("force", false, "force apply settings. Do not wait for user confirmation")
	osRhelSetCmd.Flags().StringP("attr", "A", "", "Set individual attributes for mssql-conf(ex: -A \"memory.memorylimitmb=8192\") see help for more info.")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//osRhelSetCmd.MarkFlagRequired("ip")
	//osRhelSetCmd.MarkFlagRequired("user")
	//osRhelSetCmd.MarkFlagRequired("pass")
	osRhelSetCmd.MarkFlagRequired("workload")
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

func setHugePages(client *ssh.Client, cmd string) error {
	_, err := util.ExecCmd(client, cmd)
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

func setTunedAdmSettingsOrcl(client *ssh.Client) {
	cmd1 := `
	mkdir -p /usr/lib/tuned/dell_oracle_oltp
	`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())

	cmd2 := `
cat > /usr/lib/tuned/dell_oracle_oltp/tuned.conf << "EOF"
#
# Dell Custom tuned configuration for OLTP workload on Oracle database
#

[main]
summary=Optimize for OLTP workload on Oracle RDBMS
include=throughput-performance

[sysctl]
vm.swappiness = 10
vm.dirty_background_ratio = 3
vm.dirty_ratio = 40
vm.dirty_expire_centisecs = 500
vm.dirty_writeback_centisecs = 100
kernel.shmmax = 4398046511104
kernel.shmall = 1073741824
kernel.shmmni = 4096
kernel.sem = 250 32000 100 128
fs.file-max = 6815744
fs.aio-max-nr = 1048576
net.ipv4.ip_local_port_range = 9000 65500
net.core.rmem_default = 262144
net.core.rmem_max = 4194304
net.core.wmem_default = 262144
net.core.wmem_max = 1048576
kernel.panic_on_oops = 1
kernel.numa_balancing=0

[vm]
transparent_hugepages=never
EOF
`
	_, err = execCmd(client, cmd2)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res2.String())

	cmd3 := `
chmod +x /usr/lib/tuned/dell_oracle_oltp/tuned.conf
tuned-adm profile dell_oracle_oltp
tuned-adm active
`
	res3, err := execCmd(client, cmd3)
	if err != nil {
		panic(err)
	}
	fmt.Println(res3.String())
}

func setTunedAdmSettingsSql(client *ssh.Client) {
	cmd1 := `
	mkdir -p /usr/lib/tuned/dell_mssql_oltp
	`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())

	cmd2 := `
cat > /usr/lib/tuned/dell_mssql_oltp/tuned.conf << "EOF"
#
# A Tuned configuration for OLTP workload on SQL Server running on Linux
#

[main]
summary=Optimize for OLTP workload on Microsoft SQL Server
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
chmod +x /usr/lib/tuned/dell_mssql_oltp/tuned.conf
tuned-adm profile dell_mssql_oltp
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
	fmt.Println("mssql-conf changes applied successfully... ")
}

func setMssqlDiskSettings(client *ssh.Client) {
	cmd1 := `
files=( $(find / \( -name "*.ldf" -o -name "*.mdf" -o -name "*.ndf"  \) -type f -print0 2>/dev/null |xargs -0))
declare -A allDvc
for (( i=0; i<${#files[@]}; i++ )); 
do 
	fileName=${files[i]} 
	dev=$(df $fileName | awk '/^\/dev/ {print $1}')
	if [[ $dev != "/dev/mapper/rhel-root" ]]; then
		allDvs=$dev;
	fi
done
declare -A uniqDvs
for dvs in "${allDvs[@]}"; do
	uniqDvs[$dvs]=0 
done
for dv in "${!uniqDvs[@]}"; do
uid=$(blkid ${dv} | awk '{print $2}'|sed 's/"//g')
fileData=$(grep -hnr "$uid" /etc/fstab | awk '{ print  $4  }')
NEW="defaults,noatime"
sed -i "/^$uid/s/$fileData/$NEW/g" /etc/fstab
blockdev --setra 4096 $dv
done	
	`
	_, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Disk best practices applied successfully... ")
}

func getHugepagesRecomendValue(client *ssh.Client) string {
	var err error
	hugepagesCmd := `
cat > /tmp/hugepagesRecomend.sh << "EOF"
#!/bin/bash
KERN=$(uname -r | awk -F. '{ printf("%d.%d\n",$1,$2); }')
HPG_SZ=$(grep Hugepagesize /proc/meminfo | awk {'print $2'})
NUM_PG=1
for SEG_BYTES in $(ipcs -m | awk {'print $5'} | grep "[0-9][0-9]*")
do
MIN_PG=$(echo "$SEG_BYTES/($HPG_SZ*1024)" | bc -q)
if [ $MIN_PG -gt 0 ]; then
NUM_PG=$(echo "$NUM_PG+$MIN_PG+1" | bc -q)
fi
done
case $KERN in
'2.4') HUGETLB_POOL=$(echo "$NUM_PG*$HPG_SZ/1024" | bc -q);
echo "vm.hugetlb_pool = $HUGETLB_POOL" ;;
'2.6') echo "vm.nr_hugepages = $NUM_PG" ;;
'3.'*) echo "vm.nr_hugepages = $NUM_PG" ;;
'4.'*) echo "vm.nr_hugepages = $NUM_PG" ;;
*) echo "Unrecognized kernel version $KERN. Exiting." ;;
esac
EOF
`
	_, err = execCmd(client, hugepagesCmd)
	if err != nil {
		panic(err)
	}

	cmd2 := `
chmod +x /tmp/hugepagesRecomend.sh
a=$(/tmp/hugepagesRecomend.sh | awk {'print $3'}) && echo $a
`
	res3, err := execCmd(client, cmd2)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res3.String())

	cmd3 := `
rm -f /tmp/hugepagesRecomend.sh 
`
	_, err = execCmd(client, cmd3)
	if err != nil {
		panic(err)
	}
	return res3.String()
}
