package cmd

import (
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

var tunedAdm bool = false
var mssqlConf bool = false

//var workload string = "sql"

var osRhelGetCmd = &cobra.Command{
	Use:   "get",
	Short: "This get command will validate best practice settings for RHEL operating system",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := initOsRhelGetStep(cmd, args)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		outFormat, _ := cmd.Flags().GetString("out")
		disk, _ := cmd.Flags().GetBool("disk")
		hpage, _ := cmd.Flags().GetBool("hpage")
		bps, _ := cmd.Flags().GetBool("bps")

		if workload == "sql" {
			if tunedAdm || bps {
				fmt.Println("Tuned-Adm Settings:")
				out1 := getTunedAdmSettingsSql(c)
				if outFormat == "table" {
					olog.Print(out1)
				} else if outFormat == "json" {
					fmt.Println(util.PrettyPrint(out1))
				} else if outFormat == "csv" {
					of1 := util.GetFilenameDate("osTunedAdmSettingReport", "csv")
					b1, err := csvutil.Marshal(out1)
					if err != nil {
						fmt.Println("error:", err)
					}
					util.WriteCsvReport(of1, string(b1))
				}
			}
			if disk || bps {
				fmt.Println("MSSQL Disk related Settings:")
				out3 := getMssqlDiskSettings(c)
				if outFormat == "table" {
					olog.Print(out3)
				} else if outFormat == "json" {
					fmt.Println(util.PrettyPrint(out3))
				} else if outFormat == "csv" {
					of3 := util.GetFilenameDate("os_rhel_mssql_disk_bps_Report", "csv")
					b3, err := csvutil.Marshal(out3)
					if err != nil {
						fmt.Println("error:", err)
					}
					util.WriteCsvReport(of3, string(b3))
				}
			}
			if mssqlConf || bps {
				fmt.Println("MSSQL-Conf Settings:")
				out2 := getMssqlConfSettings(c)
				if outFormat == "table" {
					olog.Print(out2)
				} else if outFormat == "json" {
					fmt.Println(util.PrettyPrint(out2))
				} else if outFormat == "csv" {
					of2 := util.GetFilenameDate("mssqlConfigReport", "csv")
					b2, err := csvutil.Marshal(out2)
					if err != nil {
						fmt.Println("error:", err)
					}
					util.WriteCsvReport(of2, string(b2))
				}
			}
			if !tunedAdm && !mssqlConf && !disk && !bps {
				fmt.Println("no sub flag (--bps or --tunedadm or --msconf or --disk) provided")
				fmt.Println("use below instruction to see help and examples for this command")
				fmt.Println("dso os rhel get -w sql --help")
			}
		}
		if workload == "oracle" {
			if tunedAdm || bps {
				fmt.Println("Current Tuned-Adm Settings:")
				out3 := getTunedAdmSettingsOrcl(c)
				if outFormat == "table" {
					olog.Print(out3)
				} else if outFormat == "json" {
					fmt.Println(util.PrettyPrint(out3))
				} else if outFormat == "csv" {
					of3 := util.GetFilenameDate("osTunedAdmSettingReport", "csv")
					b3, err := csvutil.Marshal(out3)
					if err != nil {
						fmt.Println("error:", err)
					}
					util.WriteCsvReport(of3, string(b3))
				}
			}
			if hpage || bps {
				fmt.Println("Hugepage Settings:")
				out5 := getHugePageDetails(c)
				if outFormat == "table" {
					olog.Print(out5)
				} else if outFormat == "json" {
					fmt.Println(util.PrettyPrint(out5))
				} else if outFormat == "csv" {
					of5 := util.GetFilenameDate("osHugepageSettingReport", "csv")
					b5, err := csvutil.Marshal(out5)
					if err != nil {
						fmt.Println("error:", err)
					}
					util.WriteCsvReport(of5, string(b5))
				}
			}
		}
	},
}

func init() {
	rhelCmd.AddCommand(osRhelGetCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	osRhelGetCmd.Flags().StringP("ip", "I", "", "IP / FQDN of the RHEL operating system")
	osRhelGetCmd.Flags().StringP("portSSH", "p", "22", "SSH port for connecting to RHEL os")
	osRhelGetCmd.Flags().StringP("user", "U", "", "Username for the RHEL operating system")
	osRhelGetCmd.Flags().StringP("pass", "P", "", "Password for the RHEL operating system")
	osRhelGetCmd.Flags().StringP("workload", "w", "", "Application workload [sql/oracle]")
	osRhelGetCmd.Flags().Bool("tunedadm", false, "Get setting values for tuned-Adm profile")
	osRhelGetCmd.Flags().Bool("msconf", false, "Get setting values for mssql-conf")
	osRhelGetCmd.Flags().Bool("disk", false, "Get disk related best practice settings for mssql-conf")
	osRhelGetCmd.Flags().Bool("hpage", false, "Get Hugepages settings for oracle workload")
	osRhelGetCmd.Flags().Bool("bps", false, "Get all best practice settings for given workload on OS layer")
	osRhelGetCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")
	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//osRhelGetCmd.MarkFlagRequired("ip")
	//osRhelGetCmd.MarkFlagRequired("user")
	//osRhelGetCmd.MarkFlagRequired("pass")
	osRhelGetCmd.MarkFlagRequired("workload")
}

type displaySettings struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
	OptimalValues string `json:"OptimalValues"`
	Diff          string `json:"Diff"`
}

func getHugePageDetails(client *ssh.Client) []displaySettings {
	var settingSlice []displaySettings
	cmd1 := "a=$(grep Hugepagesize /proc/meminfo | awk {'print $2'}) && echo $a"
	res1, err := execCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	hugePageRecomendation := getHugepagesRecomendValue(client)

	var setting displaySettings
	setting.Settings = "vm.nr_hugepages"
	setting.RunningValues = strings.Trim(res1.String(), "\n")
	setting.OptimalValues = strings.Trim(hugePageRecomendation, "\n")
	if setting.RunningValues != setting.OptimalValues {
		setting.Diff = "*"
	}
	settingSlice = append(settingSlice, setting)
	//fmt.Println(setting)
	//olog.Print(settingSlice)
	return settingSlice
}

type tunedAdmSettings struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
	OptimalValues string `json:"OptimalValues"`
	Diff          string `json:"Diff"`
}

type mssqlConfSettings struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
}

func initOsRhelGetStep(cmd *cobra.Command, args []string) (*ssh.Client, error) {

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

func getTunedAdmSettingsOrcl(client *ssh.Client) []tunedAdmSettings {
	cmd1 := `
	sysctl vm.swappiness vm.dirty_background_ratio vm.dirty_ratio vm.dirty_expire_centisecs vm.dirty_writeback_centisecs kernel.shmmax kernel.shmall kernel.shmmni fs.file-max fs.aio-max-nr net.core.rmem_default net.core.rmem_max net.core.wmem_default net.core.wmem_max kernel.panic_on_oops kernel.numa_balancing | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3 "\"}" }' && sysctl net.ipv4.ip_local_port_range | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3" "$4 "\"}" }' && sysctl kernel.sem | awk ' { print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3" "$4" "$5" "$6 "\"}" }'
	`
	res, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []tunedAdmSettings
	for _, m := range strings.Split(res.String(), "\n") {
		var osd tunedAdmSettings
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}

	for k := range osdata {
		//oss := osdata[k].Settings
		switch osdata[k].Settings {
		case "kernel.numa_balancing":
			osdata[k].OptimalValues = "0"
		case "kernel.shmmax":
			osdata[k].OptimalValues = "4398046511104"
		case "kernel.shmall":
			osdata[k].OptimalValues = "1073741824"
		case "kernel.shmmni":
			osdata[k].OptimalValues = "4096"
		case "kernel.sem":
			osdata[k].OptimalValues = "250 32000 100 128"
		case "net.core.rmem_default":
			osdata[k].OptimalValues = "262144"
		case "net.core.rmem_max":
			osdata[k].OptimalValues = "4194304"
		case "net.core.wmem_default":
			osdata[k].OptimalValues = "262144"
		case "net.core.wmem_max":
			osdata[k].OptimalValues = "1048576"
		case "net.ipv4.ip_local_port_range":
			osdata[k].OptimalValues = "9000 65499"
		case "vm.dirty_background_ratio":
			osdata[k].OptimalValues = "3"
		case "vm.dirty_expire_centisecs":
			osdata[k].OptimalValues = "500"
		case "vm.dirty_ratio":
			osdata[k].OptimalValues = "40"
		case "vm.dirty_writeback_centisecs":
			osdata[k].OptimalValues = "100"
		case "vm.swappiness":
			osdata[k].OptimalValues = "10"
		case "kernel.panic_on_oops":
			osdata[k].OptimalValues = "1"
		case "fs.file-max":
			osdata[k].OptimalValues = "6815744"
		case "fs.aio-max-nr":
			osdata[k].OptimalValues = "1048576"
		}
		if osdata[k].RunningValues != osdata[k].OptimalValues {
			osdata[k].Diff = "*"
		}
	}
	//olog.Print(osdata)
	return osdata
}

func getTunedAdmSettingsSql(client *ssh.Client) []tunedAdmSettings {
	cmd1 := `
	sysctl -a | grep -E 'vm.swappiness|vm.dirty_background_ratio|vm.dirty_ratio|vm.dirty_expire_centisecs|vm.dirty_writeback_centisecs|vm.transparent_hugepages|vm.max_map_count|net.core.rmem_default|net.core.rmem_max|net.core.wmem_default|net.core.wmem_max|kernel.numa_balancing|kernel.sched_min_granularity_ns|kernel.sched_wakeup_granularity_ns' |awk ' {
		print "{\"Settings\":\"" $1 "\",\"RunningValues\":\"" $3 "\"}"
		}'
	`

	res, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []tunedAdmSettings
	for _, m := range strings.Split(res.String(), "\n") {
		var osd tunedAdmSettings
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}

	for k := range osdata {
		//oss := osdata[k].Settings
		switch osdata[k].Settings {
		case "kernel.numa_balancing":
			osdata[k].OptimalValues = "0"
		case "kernel.numa_balancing_scan_delay_ms":
			osdata[k].OptimalValues = "1000"
		case "kernel.numa_balancing_scan_period_max_ms":
			osdata[k].OptimalValues = "60000"
		case "kernel.numa_balancing_scan_period_min_ms":
			osdata[k].OptimalValues = "1000"
		case "kernel.numa_balancing_scan_size_mb":
			osdata[k].OptimalValues = "256"
		case "net.core.rmem_default":
			osdata[k].OptimalValues = "262144"
		case "net.core.rmem_max":
			osdata[k].OptimalValues = "4194304"
		case "net.core.wmem_default":
			osdata[k].OptimalValues = "262144"
		case "net.core.wmem_max":
			osdata[k].OptimalValues = "1048576"
		case "vm.dirty_background_ratio":
			osdata[k].OptimalValues = "3"
		case "vm.dirty_expire_centisecs":
			osdata[k].OptimalValues = "500"
		case "vm.dirty_ratio":
			osdata[k].OptimalValues = "80"
		case "vm.dirty_writeback_centisecs":
			osdata[k].OptimalValues = "100"
		case "vm.max_map_count":
			osdata[k].OptimalValues = "1600000"
		case "vm.swappiness":
			osdata[k].OptimalValues = "1"
		case "kernel.sched_wakeup_granularity_ns":
			osdata[k].OptimalValues = "2000000"
		case "kernel.sched_min_granularity_ns":
			osdata[k].OptimalValues = "15000000"
		}
		if osdata[k].RunningValues != osdata[k].OptimalValues {
			osdata[k].Diff = "*"
		}
	}
	//olog.Print(osdata)
	return osdata
}

func getMssqlConfSettings(client *ssh.Client) []mssqlConfSettings {
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

type mssqlDiskSettings struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
	OptimalValues string `json:"OptimalValues"`
	Diff          string `json:"Diff"`
}

func getMssqlDiskSettings(client *ssh.Client) []mssqlDiskSettings {
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
echo $fileData | awk -v x=$dskOpt '{ print "{\"Settings\":\"" x "\",\"RunningValues\":\"" $4 "\",\"OptimalValues\":\"defaults,noatime\",\"Diff\":\"\"}" }'
blockdev --getra $dv | awk -v y=$dskReadAhead '{ print "{\"Settings\":\"" y "\",\"RunningValues\":\"" $1 "\",\"OptimalValues\":\"4096\",\"Diff\":\"\"}" }'
done
	`
	res, err := util.ExecCmd(client, cmd1)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res.String())
	var osdata []mssqlDiskSettings
	for _, m := range strings.Split(res.String(), "\n") {
		var osd mssqlDiskSettings
		if m != "" {
			json.Unmarshal([]byte(m), &osd)
			osdata = append(osdata, osd)
		}
	}
	/*
		for k := range osdata {
			fmt.Println(k)
			fmt.Println(osdata[k].RunningValues)
			fmt.Println(osdata[k].OptimalValues)
			if osdata[k].RunningValues == osdata[k].OptimalValues {
				osdata[k].Diff = "*"
			} else {
				osdata[k].Diff = ""
			}
		}
	*/
	for k := range osdata {
		if osdata[k].RunningValues != osdata[k].OptimalValues {
			osdata[k].Diff = "*"
		}
	}
	//olog.Print(osdata)
	return osdata
}
