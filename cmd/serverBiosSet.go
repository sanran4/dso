/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

//var srvSetJobId string

// serverCmd represents the server command
var serverBiosSetCmd = &cobra.Command{
	Use:   "set",
	Short: "This get command will set best practice settings for server layer of the solution",
	Long:  `This get command will set best practice settings for Intel based server layer of the solution`,
	Run: func(cmd *cobra.Command, args []string) {
		bpsFlag, _ := cmd.Flags().GetBool("bps")
		job_id, _ := cmd.Flags().GetString("jobStatus")
		testFlag, _ := cmd.Flags().GetBool("test")
		idracIP, user, pass, srvSetAttr, srvJobConfig, rbtSrv := initSrvSet(cmd, args)
		fmt.Println(srvSetAttr)
		//var job_id string
		if bpsFlag {
			srvSetBiosData(idracIP, user, pass, "WorkloadProfile", "DbOptimizedProfile")
			srvSetBiosData(idracIP, user, pass, "SysProfile", "PerfOptimized")
			srvSetBiosData(idracIP, user, pass, "ProcVirtualization", "Enabled")
			srvSetBiosData(idracIP, user, pass, "ProcX2Apic", "Enabled")
			srvSetBiosData(idracIP, user, pass, "LogicalProc", "Enabled")
			srvSetBiosData(idracIP, user, pass, "MemOpMode", "OptimizerMode")
			srvSetBiosData(idracIP, user, pass, "SerialComm", "Off")
			srvSetBiosData(idracIP, user, pass, "UsbPorts", "AllOff")
			srvSetBiosData(idracIP, user, pass, "UsbManagedPort", "Off")
			job_id := srvCreateBiosConfigJob(idracIP, user, pass)
			fmt.Println(job_id)
			time.Sleep(10 * time.Second)
			srvRebootServer(idracIP, user, pass)
			srvLoopJobStatus(idracIP, user, pass, job_id)
		} else if srvSetAttr != "" {
			attName, attValue := srvSetParseAttr(srvSetAttr)
			srvSetBiosData(idracIP, user, pass, attName, attValue)
			job_id := srvCreateBiosConfigJob(idracIP, user, pass)
			fmt.Println(job_id)
			time.Sleep(10 * time.Second)
			srvRebootServer(idracIP, user, pass)
			srvLoopJobStatus(idracIP, user, pass, job_id)
		} else if srvJobConfig {
			srvCreateBiosConfigJob(idracIP, user, pass)
		} else if job_id != "" {
			srvLoopJobStatus(idracIP, user, pass, job_id)
		} else if rbtSrv {
			srvRebootServer(idracIP, user, pass)
		} else if testFlag {
			srvLoopJobStatus(idracIP, user, pass, "JID_526269044866")
		}
	},
}

func init() {
	serverBiosCmd.AddCommand(serverBiosSetCmd)
	serverBiosSetCmd.Flags().StringP("idracIP", "I", "", "iDRAC IP of the server")
	serverBiosSetCmd.Flags().StringP("user", "U", "", "Username for the server iDRAC")
	serverBiosSetCmd.Flags().StringP("pass", "P", "", "Password for the server iDRAC")
	//serverSetCmd.Flags().Bool("bios", false, "work with bios section of server iDRAC")
	serverBiosSetCmd.Flags().Bool("bps", false, "set all best practice on server iDRAC")
	serverBiosSetCmd.Flags().StringP("attr", "A", "", "Attribute/s to be set on server layer")
	serverBiosSetCmd.Flags().Bool("jobcfg", false, "Work with bios job cofig server iDRAC")
	serverBiosSetCmd.Flags().StringP("jobStatus", "j", "", "Check BIOS Job Status in loop based on job_id")
	serverBiosSetCmd.Flags().Bool("reboot", false, "Reboot server using iDRAC")
	serverBiosSetCmd.Flags().Bool("test", false, "test individual function")

}

func initSrvSet(cmd *cobra.Command, args []string) (ip, u, p, a string, jc, rb bool) {
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
	//srvSetBios, _ = cmd.Flags().GetBool("bios")
	srvJobConfig, _ := cmd.Flags().GetBool("jobcfg")
	srvSetAttr, _ := cmd.Flags().GetString("attr")
	rbtSrv, _ := cmd.Flags().GetBool("reboot")

	return idracIP, user, pass, srvSetAttr, srvJobConfig, rbtSrv
}

func srvSetParseAttr(str string) (attr, val string) {
	tmp := strings.Split(str, ":")
	attr = tmp[0]
	val = tmp[1]
	return
}

func srvSetBiosData(idracIP, user, pass, aName, aValue string) {
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
	baseURL := "https://" + idracIP + "/redfish/v1/Systems/System.Embedded.1/Bios/Settings"
	var jsonStr = []byte("{\"Attributes\":{\"" + aName + "\":\"" + aValue + "\"}}")
	req, err := http.NewRequest("PATCH", baseURL, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Printf("error creating PATCH request %v", err)
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
	if resp.StatusCode == 200 {
		fmt.Printf("Command passed to set BIOS attribute %s pending value to %s\n", aName, aValue)
	} else {
		fmt.Printf("Command failed, errror is %s\n", resp.Status)
		fmt.Println(string(responseBytes))
	}
	//fmt.Println(resp.Body)
}

// Function to create BIOS target config job
func srvCreateBiosConfigJob(idracIP, user, pass string) (srvSetJobId string) {
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

	jobUrl := "https://" + idracIP + "/redfish/v1/Managers/iDRAC.Embedded.1/Jobs"
	//payload := {"TargetSettingsURI":"/redfish/v1/Systems/System.Embedded.1/Bios/Settings"}
	var payload = []byte("{\"TargetSettingsURI\":\"/redfish/v1/Systems/System.Embedded.1/Bios/Settings\"}")
	req, err := http.NewRequest("POST", jobUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("error creating PATCH request %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error requesting data %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Header.Get("Location"))
	job_loc := resp.Header.Get("Location")
	job_array := strings.Split(job_loc, "/")
	srvSetJobId = job_array[len(job_array)-1]

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	if resp.StatusCode == 200 {
		fmt.Printf("Command passed to create target config job, status code 200 returned\n")
	} else {
		fmt.Printf("Command failed, errror is %s\n", resp.Status)
		fmt.Println(string(responseBytes))
	}
	//fmt.Println(string(responseBytes))
	return srvSetJobId
}

type srvJobStatus struct {
	Id       string `json:"Id"`
	JobType  string `json:"JobType"`
	JobState string `json:"JobState"`
	Message  string `json:"Message"`
}

// Function to check job status
func srvLoopJobStatus(idracIP, user, pass, jobId string) {

	iterateCnt := 30

	for iterateCnt > 0 {
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

		jobUrl := "https://" + idracIP + "/redfish/v1/Managers/iDRAC.Embedded.1/Jobs/" + jobId
		req, err := http.NewRequest("GET", jobUrl, nil)
		if err != nil {
			log.Printf("error creating GET request %v", err)
			os.Exit(1)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json;charset=utf-8")
		req.Header.Add("Cache-Control", "no-cache")
		req.SetBasicAuth(user, pass)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error requesting data %v", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Could not read response body. %v", err)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			fmt.Printf("Failed to get job status, errror is %s\n", resp.Status)
			fmt.Println(string(rb))
			os.Exit(1)
		}
		var sjs srvJobStatus
		if err := json.Unmarshal(rb, &sjs); err != nil {
			log.Printf("Could not unmarshal srvJobStatus. %v", err)
		}

		jobStaus := sjs.JobState
		if jobStaus == "Failed" || jobStaus == "Completed" {
			fmt.Printf("Job Status %s . Exiting execution\n", jobStaus)
			break
		}

		fmt.Printf("Job Status %s . Continue to check for operation\n", jobStaus)
		time.Sleep(30 * time.Second)
		iterateCnt -= 1
	}
	//fmt.Println(string(responseBytes))
}
