package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

var serverRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "This reboot command will reboot physical server within solution using iDRAC",
	Example: `
EX1: dso server reboot -I 10.0.0.1 -U user1 -P pass1
`,
	Run: func(cmd *cobra.Command, args []string) {
		initSrvReboot(cmd, args)
	},
}

func init() {
	serverCmd.AddCommand(serverRebootCmd)
	serverRebootCmd.Flags().StringP("idracIP", "I", "", "iDRAC IP of the server")
	serverRebootCmd.Flags().StringP("user", "U", "", "Username for the server iDRAC")
	serverRebootCmd.Flags().StringP("pass", "P", "", "Password for the server iDRAC")
}

func initSrvReboot(cmd *cobra.Command, args []string) {
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

	srvRebootServer(idracIP, user, pass)
}

// Function to Reboot Server through iDRAC
func srvRebootServer(idracIP, user, pass string) {
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

	jobUrl := "https://" + idracIP + "/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset"
	payload := []byte("{\"ResetType\":\"ForceOff\"}")
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
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	if resp.StatusCode == 204 {
		fmt.Printf("Command passed to power OFF server")
	} else {
		fmt.Printf("Command failed, errror is %s\n", resp.Status)
		fmt.Println(string(responseBytes))
	}
	fmt.Println(string(responseBytes))

	time.Sleep(10 * time.Second)

	payload = []byte("{\"ResetType\":\"On\"}")
	req2, err := http.NewRequest("POST", jobUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("error creating PATCH request %v", err)
	}
	req2.Header.Add("Content-Type", "application/json")
	req2.Header.Add("Accept", "application/json;charset=utf-8")
	req2.Header.Add("Cache-Control", "no-cache")
	req2.SetBasicAuth(user, pass)
	resp2, err := client.Do(req2)
	if err != nil {
		log.Printf("error requesting data %v", err)
	}
	defer resp2.Body.Close()
	rb, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	if resp2.StatusCode == 204 {
		fmt.Printf("Command passed to power ON server")
	} else {
		fmt.Printf("Command failed, errror is %s\n", resp2.Status)
		fmt.Println(string(rb))
	}
	fmt.Println(string(rb))
}
