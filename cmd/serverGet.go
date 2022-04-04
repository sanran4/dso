/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/da0x/golang/olog"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

var srvGetBios bool = false

// serverCmd represents the server command
var serverGetCmd = &cobra.Command{
	Use:   "get",
	Short: "This get command will pull best practice settings for server layer of the solution",
	Long:  `This get command will pull best practice settings for Intel based server layer of the solution`,
	Example: `
EX1: dso server get --bios -I 10.0.0.1 -U user1 -P pass1
EX2: dso server get --bios -I 10.0.0.1 --bios -U user1 -P pass1
EX3: dso server get --bios --idracIP=10.0.0.1 --user=user1 --pass=pass1
`,
	Run: func(cmd *cobra.Command, args []string) {
		initSrvGet(cmd, args)
	},
}

func init() {
	serverCmd.AddCommand(serverGetCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	serverGetCmd.Flags().StringP("idracIP", "I", "", "iDRAC IP of the server")
	serverGetCmd.Flags().StringP("user", "U", "", "Username for the server iDRAC")
	serverGetCmd.Flags().StringP("pass", "P", "", "Password for the server iDRAC")
	serverGetCmd.Flags().Bool("bios", false, "Get setting values for tuned-Adm profile")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//serverGetCmd.MarkFlagRequired("idracIP")
	//serverGetCmd.MarkFlagRequired("user")
	//reportCmd.MarkFlagRequired("pass")
}

type srvGetBiosSetting struct {
	Settings      string `json:"Settings"`
	RunningValues string `json:"RunningValues"`
	OptimalValues string `json:"OptimalValues"`
	Diff          string `json:"Diff"`
}

type BiosAttribute struct {
	MemOpMode          string `json:"MemOpMode"`
	LogicalProc        string `json:"LogicalProc"`
	ProcVirtualization string `json:"ProcVirtualization"`
	ProcX2Apic         string `json:"ProcX2Apic"`
	SysProfile         string `json:"SysProfile"`
	WorkloadProfile    string `json:"WorkloadProfile"`
}
type srvGetBiosConfig struct {
	Name       string        `json:"name"`
	Attributes BiosAttribute `json:"Attributes"`
}

func srvGetBiosData(baseURL, user, pass string) {

	// TODO: This is insecure; use only in dev environments.
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	//client := &http.Client{Transport: tr}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Printf("error creating GET request %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Printf("error requesting data %v", err)
	}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	var srvgb srvGetBiosConfig
	if err := json.Unmarshal(responseBytes, &srvgb); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	//fmt.Println(PrettyPrint(biosConfig))
	attr := BiosAttribute{}
	attr = srvgb.Attributes

	var sgbs []srvGetBiosSetting

	fields := reflect.TypeOf(attr)
	values := reflect.ValueOf(attr)
	num := fields.NumField()
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		var l_sgbs srvGetBiosSetting

		var v string
		switch value.Kind() {
		case reflect.String:
			v = value.String()
		case reflect.Int:
			v = (strconv.FormatInt(value.Int(), 10))
		case reflect.Int32:
			v = strconv.FormatInt(value.Int(), 10)
		case reflect.Int64:
			v = strconv.FormatInt(value.Int(), 10)
		default:
			v = value.String()
		}
		l_sgbs.Settings = field.Name
		l_sgbs.RunningValues = v
		switch field.Name {
		case "MemOpMode":
			l_sgbs.OptimalValues = "OptimizerMode"
		case "LogicalProc":
			l_sgbs.OptimalValues = "Enabled"
		case "ProcVirtualization":
			l_sgbs.OptimalValues = "Enabled"
		case "ProcX2Apic":
			l_sgbs.OptimalValues = "Enabled"
		case "SysProfile":
			l_sgbs.OptimalValues = "Performance"
		case "WorkloadProfile":
			l_sgbs.OptimalValues = "Database Optimized Performance Profile"
		}
		//fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		sgbs = append(sgbs, l_sgbs)
	}

	for k := range sgbs {
		if sgbs[k].RunningValues != sgbs[k].OptimalValues {
			sgbs[k].Diff = "*"
		}
	}
	//return attr
	olog.Print(sgbs)
}

func initSrvGet(cmd *cobra.Command, args []string) {
	idracIP, ok := os.LookupEnv("SERVER_IDRAC_HOST")
	if !ok {
		idracIP, _ = cmd.Flags().GetString("idracIP")
	}
	user, ok := os.LookupEnv("SERVER_IDRAC_USER")
	if !ok {
		user, _ = cmd.Flags().GetString("user")
	}
	//idracIP, _ := cmd.Flags().GetString("idracIP")
	//user, _ := cmd.Flags().GetString("user")
	pass, _ := cmd.Flags().GetString("pass")
	var err error
	if pass == "" {
		pass, err = util.GetPasswd()
		if err != nil {
			log.Printf("error getting password %v", err)
		}
	}
	srvGetBios, _ = cmd.Flags().GetBool("bios")
	var baseURL string
	if srvGetBios {
		baseURL = "https://" + idracIP + "/redfish/v1/Systems/System.Embedded.1/Bios"
		srvGetBiosData(baseURL, user, pass)
	}
}
