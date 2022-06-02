package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/da0x/golang/olog"
	"github.com/jszwec/csvutil"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var osRhelReportCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will pull report from the RHEL operating system",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := InitialSetup(cmd, args)
		if err != nil {
			panic(err)
		}
		defer client.Close()
		outFormat, _ := cmd.Flags().GetString("out")
		if workload == "sql" {
			output1 := getSysctlConfigSql(client)
			output2 := getMssqlConfSettingsReport(client)
			output3 := getMssqlDiskSettingsReport(client)
			if outFormat == "table" {
				fmt.Println("OS kernal settings")
				olog.Print(output1)
				fmt.Println("Disk settings")
				olog.Print(output3)
				fmt.Println("SQL Server config settings")
				olog.Print(output2)
			} else if outFormat == "json" {
				fmt.Println("OS kernal settings")
				fmt.Println(util.PrettyPrint(output1))
				fmt.Println("Disk settings")
				fmt.Println(util.PrettyPrint(output3))
				fmt.Println("SQL Server config settings")
				fmt.Println(util.PrettyPrint(output2))
			} else if outFormat == "csv" {
				outputFile1 := util.GetFilenameDate("osKernalSettingReport", "csv")
				b1, err := csvutil.Marshal(output1)
				if err != nil {
					fmt.Println("error:", err)
				}
				util.WriteCsvReport(outputFile1, string(b1))
				outputFile3 := util.GetFilenameDate("sqlServerDiskSettingReport", "csv")
				b3, err := csvutil.Marshal(output3)
				if err != nil {
					fmt.Println("error:", err)
				}
				util.WriteCsvReport(outputFile3, string(b3))
				outputFile2 := util.GetFilenameDate("sqlServerConfigReport", "csv")
				b2, err := csvutil.Marshal(output2)
				if err != nil {
					fmt.Println("error:", err)
				}
				util.WriteCsvReport(outputFile2, string(b2))
			}

		}
		if workload == "oracle" {
			output1 := getSysctlConfigOrcl(client)
			output2 := getHugePageDetailsReport(client)
			if outFormat == "table" {
				fmt.Println("OS kernal settings")
				olog.Print(output1)
				fmt.Println("HugePages settings")
				olog.Print(output2)
			} else if outFormat == "json" {
				fmt.Println("OS kernal settings")
				fmt.Println(util.PrettyPrint(output1))
				fmt.Println("HugePages settings")
				fmt.Println(util.PrettyPrint(output2))
			} else if outFormat == "csv" {
				outputFile1 := util.GetFilenameDate("osKernalSettingReport", "csv")
				b1, err := csvutil.Marshal(output1)
				if err != nil {
					fmt.Println("error:", err)
				}
				util.WriteCsvReport(outputFile1, string(b1))
				outputFile2 := util.GetFilenameDate("hugePagesSettingReport", "csv")
				b2, err := csvutil.Marshal(output2)
				if err != nil {
					fmt.Println("error:", err)
				}
				util.WriteCsvReport(outputFile2, string(b2))
			}
		}

	},
}

func init() {
	rhelCmd.AddCommand(osRhelReportCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	osRhelReportCmd.Flags().StringP("ip", "I", "", "IP / FQDN of the RHEL operating system")
	osRhelReportCmd.Flags().StringP("portSSH", "p", "22", "SSH port for connecting to RHEL os")
	osRhelReportCmd.Flags().StringP("user", "U", "", "Username for the RHEL operating system")
	osRhelReportCmd.Flags().StringP("pass", "P", "", "Password for the RHEL operating system")
	osRhelReportCmd.Flags().StringP("workload", "w", "", "Application workload [sql/orcl]")
	osRhelReportCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")
	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//osRhelReportCmd.MarkFlagRequired("ip")
	//osRhelReportCmd.MarkFlagRequired("user")
	//osRhelReportCmd.MarkFlagRequired("pass")
	osRhelReportCmd.MarkFlagRequired("workload")
}

func execCmd(client *ssh.Client, query string) (bytes.Buffer, error) {

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(query); err != nil {
		log.Fatal("Failed to run: " + err.Error())
		return b, err
	}
	//fmt.Println(b.String())
	return b, nil
}

func InitialSetup(cmd *cobra.Command, args []string) (*ssh.Client, error) {

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

type OsSetting struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
}

func getSysctlConfigOrcl(client *ssh.Client) []OsSetting {
	cmd1 := `
	sysctl vm.swappiness vm.dirty_background_ratio vm.dirty_ratio vm.dirty_expire_centisecs vm.dirty_writeback_centisecs kernel.shmmax kernel.shmall kernel.shmmni fs.file-max fs.aio-max-nr net.core.rmem_default net.core.rmem_max net.core.wmem_default net.core.wmem_max kernel.panic_on_oops kernel.numa_balancing | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3 "\"}" }' && sysctl net.ipv4.ip_local_port_range | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3" "$4 "\"}" }' && sysctl kernel.sem | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3" "$4" "$5" "$6 "\"}" }'
	`
	res, err := execCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []OsSetting
	for _, m := range strings.Split(res.String(), "\n") {
		var osd OsSetting
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}

	//olog.Print(osdata)
	return osdata
}

func getSysctlConfigSql(client *ssh.Client) []OsSetting {
	cmd1 := `
	sysctl -a | grep -E 'vm.swappiness|vm.dirty_background_ratio|vm.dirty_ratio|vm.dirty_expire_centisecs|vm.dirty_writeback_centisecs|vm.transparent_hugepages|vm.max_map_count|net.core.rmem_default|net.core.rmem_max|net.core.wmem_default|net.core.wmem_max|kernel.numa_balancing|kernel.sched_min_granularity_ns|kernel.sched_wakeup_granularity_ns' |awk ' {
		print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3 "\"}"
		}'
	`

	res, err := execCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []OsSetting
	for _, m := range strings.Split(res.String(), "\n") {
		var osd OsSetting
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}

	//olog.Print(osdata)
	return osdata
}

func getMssqlConfSettingsReport(client *ssh.Client) []mssqlConfSettings {
	cmd1 := `
{ if echo "$(/opt/mssql/bin/mssql-conf get EULA)" | grep -q "No"; 
then echo "accepteula : NotSet"; 
else echo "$(/opt/mssql/bin/mssql-conf get EULA)"; 
fi;
if echo "$(/opt/mssql/bin/mssql-conf get control)" | grep -q "No"; 
then echo "control : NotSet"; 
else echo "$(/opt/mssql/bin/mssql-conf get control)"; 
fi;
if echo "$(/opt/mssql/bin/mssql-conf get memory)" | grep -q "No"; 
then echo "memory : NotSet"; 
else echo "$(/opt/mssql/bin/mssql-conf get memory)"; 
fi; 
if echo "$(/opt/mssql/bin/mssql-conf get traceflag)" | grep -q "No"; 
then echo "traceflag : NotSet"; 
else echo "$(/opt/mssql/bin/mssql-conf get traceflag)"; 
fi; } | awk ' {
print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3 "\"}"
}'
	`

	res, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []mssqlConfSettings
	for _, m := range strings.Split(res.String(), "\n") {
		var osd mssqlConfSettings
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}

	//olog.Print(osdata)
	return osdata
}

func getMssqlDiskSettingsReport(client *ssh.Client) []OsSetting {
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
dskOpt=$dv"_diskMountOption"
dskReadAhead=$dv"_diskReadAheadValue"
uid=$(blkid ${dv} | awk '{print $2}'|sed 's/"//g')
fileData=$(grep -hnr "$uid" /etc/fstab)
echo $fileData | awk -v x=$dskOpt '{ print "{\"Settings\":\"" x "\",\"RunningValues\":\"" $4 "\"}" }'
blockdev --getra $dv | awk -v y=$dskReadAhead '{ print "{\"Settings\":\"" y "\",\"RunningValues\":\"" $1 "\"}" }'
done
	`
	res, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []OsSetting
	for _, m := range strings.Split(res.String(), "\n") {
		var osd OsSetting
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}
	//olog.Print(osdata)
	return osdata
}

func getHugePageDetailsReport(client *ssh.Client) []OsSetting {
	var settingSlice []OsSetting
	cmd1 := "a=$(grep Hugepagesize /proc/meminfo | awk {'print $2'}) && echo $a"
	res1, err := execCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	var setting OsSetting
	setting.Settings = "vm.nr_hugepages"
	setting.RunningValues = strings.Trim(res1.String(), "\n")
	settingSlice = append(settingSlice, setting)
	//fmt.Println(setting)
	//olog.Print(settingSlice)
	return settingSlice
}
