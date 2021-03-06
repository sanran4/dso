/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/da0x/golang/olog"
	"github.com/jszwec/csvutil"
	"github.com/sanran4/dso/util"

	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var serverBiosReportCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will pull report for physical server BIOS",
	Run: func(cmd *cobra.Command, args []string) {
		outFormat, _ := cmd.Flags().GetString("out")
		output := ShowBios(cmd, args)
		if outFormat == "table" {
			olog.Print(output)
		} else if outFormat == "json" {
			fmt.Println(util.PrettyPrint(output))
		} else if outFormat == "csv" {
			outputFile := util.GetFilenameDate("serverReport", "csv")
			b, err := csvutil.Marshal(output)
			if err != nil {
				fmt.Println("error:", err)
			}
			util.WriteCsvReport(outputFile, string(b))
		}
	},
}

func init() {
	serverBiosCmd.AddCommand(serverBiosReportCmd)

	// Flags
	// Format: biosCmd.PersistentFlags().StringP(name string, shorthand string, value string, usage string)
	serverBiosReportCmd.Flags().StringP("idracIP", "I", "", "iDRAC IP of the server")
	serverBiosReportCmd.Flags().StringP("user", "U", "", "Username for the server iDRAC")
	serverBiosReportCmd.Flags().StringP("pass", "P", "", "Password for the server iDRAC")
	serverBiosReportCmd.Flags().StringP("out", "o", "table", "output format, available options (json, [table], csv)")

	//birthdayCmd.PersistentFlags().StringP("alertType", "y", "", "Possible values: email, sms")
	// Making Flags Required
	//serverReportCmd.MarkFlagRequired("idracIP")
	//serverReportCmd.MarkFlagRequired("user")
	//reportCmd.MarkFlagRequired("pass")
}

type AttrConv struct {
	ServerSetting string
	Value         string
}

type Attribute struct {
	SystemModelName     string `json:"SystemModelName" csv:"SystemModelName"`
	SystemBiosVersion   string `json:"SystemBiosVersion" csv:"SystemBiosVersion"`
	SystemServiceTag    string `json:"SystemServiceTag" csv:"SystemServiceTag"`
	SystemManufacturer  string `json:"SystemManufacturer" csv:"SystemManufacturer"`
	SysMfrContactInfo   string `json:"SysMfrContactInfo" csv:"SysMfrContactInfo"`
	SystemCpldVersion   string `json:"SystemCpldVersion" csv:"SystemCpldVersion"`
	SysMemSize          string `json:"SysMemSize" csv:"SysMemSize"`
	SysMemType          string `json:"SysMemType" csv:"SysMemType"`
	SysMemSpeed         string `json:"SysMemSpeed" csv:"SysMemSpeed"`
	MemOpMode           string `json:"MemOpMode" csv:"MemOpMode"`
	ProcBrand           string `json:"Proc1Brand" csv:"Proc1Brand"`
	ProcCoreSpeed       string `json:"ProcCoreSpeed" csv:"ProcCoreSpeed"`
	ProcNumCores        int    `json:"Proc1NumCores" csv:"Proc1NumCores"`
	ProcBusSpeed        string `json:"ProcBusSpeed" csv:"ProcBusSpeed"`
	LogicalProc         string `json:"LogicalProc" csv:"LogicalProc"`
	ProcVirtualization  string `json:"ProcVirtualization" csv:"ProcVirtualization"`
	ProcX2Apic          string `json:"ProcX2Apic" csv:"ProcX2Apic"`
	ControlledTurbo     string `json:"ControlledTurbo" csv:"ControlledTurbo"`
	NvmeMode            string `json:"NvmeMode" csv:"NvmeMode"`
	BootMode            string `json:"BootMode" csv:"BootMode"`
	SysProfile          string `json:"SysProfile" csv:"SysProfile"`
	WorkloadProfile     string `json:"WorkloadProfile" csv:"WorkloadProfile"`
	SecureBoot          string `json:"SecureBoot" csv:"SecureBoot"`
	SerialCommunication string `json:"SerialComm" csv:"SerialComm"`
	UsbPorts            string `json:"UsbPorts" csv:"UsbPorts"`
	UsbManagedPort      string `json:"UsbManagedPort" csv:"UsbManagedPort"`
}

type BiosConfig struct {
	Name       string    `json:"name"`
	Attributes Attribute `json:"Attributes"`
}

func createURL(cmd *cobra.Command, args []string) (uri, usr, pas string) {
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

	//fmt.Printf("idracIP %s\nuser %s\npass %s\n", idracIP, user, pass)
	baseURL := "https://" + idracIP + "/redfish/v1/Systems/System.Embedded.1/Bios"
	return baseURL, user, pass
}

func fetchReportData(baseURL, user, pass string) Attribute {

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	// TODO: This is insecure; use only in dev environments.
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//client := &http.Client{Transport: tr}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	//fmt.Println(baseURL)
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Printf("error creating GET request %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(user, pass)

	//req.Header = http.Header{
	//	"Content-Type": []string{"application/json;odata.metadata=minimal;charset=utf-8"},
	//}
	resp, err := client.Do(req)
	if err != nil {
		// handle err
		log.Printf("error requesting data %v", err)
	}
	//if resp.StatusCode == 401 {
	//	fmt.Printf("\nWARNING, status code %d returned. Incorrect iDRAC username/password or invalid privilege detected.", resp.StatusCode)
	//	log.Panic("status code 401")
	//}else if resp.StatusCode != 200{
	//	fmt.Printf("\nWARNING, iDRAC version installed does not support this feature using Redfish API")
	//	log.Panic("status code 401")
	//}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	//fmt.Println(string(responseBytes))

	var biosConfig BiosConfig
	if err := json.Unmarshal(responseBytes, &biosConfig); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	//fmt.Println(PrettyPrint(biosConfig))
	var attr Attribute = biosConfig.Attributes

	//fmt.Println(attr)
	return attr
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func ShowBios(cmd *cobra.Command, args []string) []AttrConv {

	//var attr []Attribute
	uri, usr, pas := createURL(cmd, args)
	rd := fetchReportData(uri, usr, pas)

	var test1 []AttrConv
	fields := reflect.TypeOf(rd)
	values := reflect.ValueOf(rd)
	num := fields.NumField()
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		var test2 AttrConv

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

		test2.ServerSetting = field.Name
		test2.Value = v
		//fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
		test1 = append(test1, test2)
	}
	//olog.Print(test1)
	return test1
}
