package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/da0x/golang/olog"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var tunedAdm bool = false
var mssqlConf bool = false

var osRhelGetCmd = &cobra.Command{
	Use:   "get",
	Short: "This get command will pull best practice settings for RHEL OS",
	Long:  `This get command will pull best practice specific settings from the RHEL operating system within your solution`,
	Example: `
Ex1: dso os rhel get --tunedadm -I 10.0.0.1 -U user1 
Ex2: dso os rhel get --tunedadm -I 10.0.0.1 -U user1 -P pass1
Ex3: dso os rhel get --msconf -I 10.0.0.1 -U user1 
Ex4: dso os rhel get --msconf --tunedadm -I 10.0.0.1 -U user1
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := initOsRhelGetStep(cmd, args)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		if tunedAdm {
			fmt.Println("Tuned-Adm Settings:")
			getTunedAdmSettings(c)
		}
		if mssqlConf {
			fmt.Println("MSSQL-Conf Settings:")
			getMssqlConfSettings(c)
		}
		if !tunedAdm && !mssqlConf {
			fmt.Println("no sub flag (--tunedadm or --msconf) provided")
			fmt.Println("use below instruction to see help and examples for this command")
			fmt.Println("dso os rhel get --help")
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
	osRhelGetCmd.Flags().Bool("tunedadm", false, "Get setting values for tuned-Adm profile")
	osRhelGetCmd.Flags().Bool("msconf", false, "Get setting values for mssql-conf")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	osRhelGetCmd.MarkFlagRequired("ip")
	osRhelGetCmd.MarkFlagRequired("user")
	//osRhelGetCmd.MarkFlagRequired("pass")
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

	ip, _ := cmd.Flags().GetString("ip")
	portSSH, _ := cmd.Flags().GetString("portSSH")
	user, _ := cmd.Flags().GetString("user")
	pass, _ := cmd.Flags().GetString("pass")
	var err error
	if pass == "" {
		pass, err = util.GetPasswd()
		if err != nil {
			log.Printf("error getting password %v", err)
		}
	}

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

func getTunedAdmSettings(client *ssh.Client) {
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
	olog.Print(osdata)
}

func getMssqlConfSettings(client *ssh.Client) {
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

	olog.Print(osdata)
}
