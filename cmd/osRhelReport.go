package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/da0x/golang/olog"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var osRhelReportCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will pull general report for the RHEL OS",
	Long:  `This report command will pull general report from the RHEL operating system within your solution`,
	Example: `
Ex1: dso os rhel report -I 10.0.0.1 -U user1 -P pass1
Ex2: dso os rhel report -I 10.0.0.1 -U user1
Ex3: dso os rhel report --ip=10.0.0.1 --user=user1 --pass=pass1
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := InitialSetup(cmd, args)
		if err != nil {
			panic(err)
		}
		defer client.Close()
		getSysctlConfig(client)
		getMssqlConfSettings(client)
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

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	osRhelReportCmd.MarkFlagRequired("ip")
	osRhelReportCmd.MarkFlagRequired("user")
	//osRhelReportCmd.MarkFlagRequired("pass")
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
	OptimalValues string `json:"OptimalValues"`
	Diff          string `json:"Diff"`
}

func getSysctlConfig(client *ssh.Client) {
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
