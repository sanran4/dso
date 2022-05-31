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
	"time"

	"github.com/da0x/golang/olog"
	"github.com/sanran4/dso/util"
	"github.com/spf13/cobra"
)

//var hgroup string

var pstoreRptCmd = &cobra.Command{
	Use:   "report",
	Short: "This report command will pull general report for PowerStore storage layer of the solution",
	Long:  `This report command will pull general report for PowerStore storage layer of the solution`,
	Run: func(cmd *cobra.Command, args []string) {
		//ShowPstoreHostGroup(cmd, args)
		//uri, usr, pas := createPstoreURL(cmd, args)
		//fetchPstoreVolRptData(uri, usr, pas)
		pstoreip, usr, pas, hgroup := parsePstoreFlag(cmd, args)
		if hgroup != "" {
			uri := buildPstoreUrl("hostGroup", pstoreip, hgroup)
			fmt.Println("PowerStore storage Host Group details")
			fetchPstoreHgroupRptData(uri, usr, pas, hgroup)
			uri = buildPstoreUrl("hostVol", pstoreip, hgroup)
			fmt.Println("PowerStore storage Host Group Volume details")
			fetchPstoreVolRptData(uri, usr, pas, hgroup)
		} else {
			uri := buildPstoreUrl("default", pstoreip, hgroup)
			fetchPstoreHgroupRptData(uri, usr, pas, hgroup)
			uri = buildPstoreUrl("hostVol", pstoreip, hgroup)
			fetchPstoreVolRptData(uri, usr, pas, hgroup)
		}

	},
}

func init() {
	pstoreCmd.AddCommand(pstoreRptCmd)

	pstoreRptCmd.Flags().StringP("ip", "I", "", "IP of the PowerStore storage")
	pstoreRptCmd.Flags().StringP("user", "U", "", "Username for PowerStore storage")
	pstoreRptCmd.Flags().StringP("pass", "P", "", "Password for PowerStore storage")
	pstoreRptCmd.Flags().StringP("hgroup", "g", "", "Host group name for PowerStore")
}

type PstoreHost struct {
	Host   string `json:"name" header:"Host"`
	OSType string `json:"os_type" header:"OSType"`
}

type displayPstoreHost struct {
	HostGroup string
	Host      string
	OSType    string
}

type PstoreHostGroup struct {
	HostGroup string       `json:"name" header:"HostGroup"`
	Hosts     []PstoreHost `json:"hosts" header:"Hosts"`
}

type PstoreHostName struct {
	HostName string `json:"name" header:"HostName"`
}

type PstoreHostData struct {
	Name string `json:"name" header:"HostGroupName"`
}

type PstoreVolume struct {
	//HostGroupName string
	VolumeName   string `json:"name" header:"VolumeName"`
	VolumeType   string `json:"type" header:"VolumeType"`
	VolumeWWN    string `json:"wwn" header:"VolumeWWN"`
	State        string `json:"state" header:"State"`
	Size         int64  `json:"size" header:"Size"`
	NodeAffinity string `json:"node_affinity" header:"NodeAffinity"`
}

type DisplayPstoreVolume struct {
	HostGroupName string
	VolumeName    string
	VolumeType    string
	VolumeWWN     string
	State         string
	Size          int64
	NodeAffinity  string
}

type PstoreHostGroupVolumeMapping struct {
	HostGroup PstoreHostData `json:"host_group" header:"HostGroup"`
	HostName  PstoreHostName `json:"host" header:"HostGroup"`
	Volume    PstoreVolume   `json:"volume" header:"Volume"`
}

func buildPstoreUrl(urlType, ip, hgroup string) string {
	//fmt.Printf("idracIP %s\nuser %s\npass %s\n", idracIP, user, pass)
	var baseURL string
	if urlType == "hostGroup" {
		baseURL = "https://" + ip + "/api/rest/host_group?select=name,hosts(name,os_type)&name=in.(" + hgroup + ")"
	} else if urlType == "hostVol" {
		baseURL = "https://" + ip + "/api/rest/host_volume_mapping?select=host_group(name),host(name),volume(name,type,wwn,state,size,node_affinity)"
	} else {
		baseURL = "https://" + ip + "/api/rest/host_group?select=name,hosts(name,os_type)"
	}
	//baseURL := "https://" + pstoreip + "/api/rest/host_group?select=name,hosts(name,os_type)"
	return baseURL
}

func parsePstoreFlag(cmd *cobra.Command, args []string) (pstoreip, usr, pas, hgroup string) {
	pstoreip, ok := os.LookupEnv("STORAGE_PSTORE_HOST")
	if !ok {
		pstoreip, _ = cmd.Flags().GetString("ip")
	}
	user, ok := os.LookupEnv("STORAGE_PSTORE_USER")
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
	hgroup, _ = cmd.Flags().GetString("hgroup")

	return pstoreip, user, pass, hgroup
}

func fetchPstoreHgroupRptData(baseURL, user, pass, hgroup string) {

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	//fmt.Println(string(baseURL))
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
		log.Printf("error requesting data %v", err)
	}

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	//fmt.Println(string(responseBytes))

	var hostGroup []PstoreHostGroup
	if err := json.Unmarshal(responseBytes, &hostGroup); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	if hgroup != "" {
		var phost []PstoreHost
		var sdph []displayPstoreHost
		//printer.Print(rd)
		//olog.Print(hostGroup)
		for _, r := range hostGroup {
			phost = append(phost, r.Hosts...)
		}
		for _, h := range phost {
			var dph displayPstoreHost
			dph.HostGroup = hgroup
			dph.Host = h.Host
			dph.OSType = h.OSType
			sdph = append(sdph, dph)
		}

		//printer.Print(phost)
		olog.Print(sdph)
	} else {
		//printer.Print(rd)
		olog.Print(hostGroup)
	}
}

func fetchPstoreVolRptData(baseURL, user, pass, hgroup string) {

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

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
		log.Printf("error requesting data %v", err)
	}

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	//fmt.Println(string(responseBytes))

	var hgvm []PstoreHostGroupVolumeMapping
	if err := json.Unmarshal(responseBytes, &hgvm); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}
	//fmt.Println(hgvm)

	var psv []PstoreVolume
	var dpsv []DisplayPstoreVolume
	//var phgn []PstoreHostData
	for _, v := range hgvm {
		var newDpsv DisplayPstoreVolume
		if hgroup != "" && v.HostGroup.Name == hgroup {
			//continue
			//psv = append(psv, v.Volume)
			newDpsv.HostGroupName = v.HostGroup.Name
			newDpsv.VolumeName = v.Volume.VolumeName
			newDpsv.VolumeType = v.Volume.VolumeType
			newDpsv.VolumeWWN = v.Volume.VolumeWWN
			newDpsv.State = v.Volume.State
			newDpsv.Size = v.Volume.Size
			newDpsv.NodeAffinity = v.Volume.NodeAffinity
			dpsv = append(dpsv, newDpsv)
		} else {
			psv = append(psv, v.Volume)
		}
		//psv = append(psv, v.Volume)
		//phgn = append(phgn, v.HostGroup)
		//} //else {
		//	psv = append(psv, v.Volume)
		//}

	}
	//olog.Print(phgn)
	//olog.Print(psv)
	if hgroup != "" {
		olog.Print(dpsv)
	} else {
		olog.Print(hgvm)
		olog.Print(psv)
	}

}

/*
func ShowPstoreHostGroup(uri, usr, pas, hgroup string) {
	//uri, usr, pas := createPstoreURL(cmd, args)
	rd := fetchPstoreHgroupRptData(uri, usr, pas)
	//frd := flat.Flatten(rd)
	//olog.Print(frd)
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	//printer.Print(frd)
	//for _, r := range rd {
	//	hs := r.Hosts
	//	for _, host := range hs {
	//		fhs := flat.Flatten(host)
	//		printer.Print(fhs)
	//	}
	//fhs := flat.Flatten(hs)
	//	o1 := flat.Flatten(r)
	//	printer.Print(o1)
	//}

	//for v := range rd {
	//	out, err := gojsonexplode.Explodejsonstr(String(k), ".")
	//	if err != nil {
	//		// handle error
	//	}
	//	fmt.Println(out)
	//}

	//fmt.Println(rd)
	//flat, err := flatten.FlattenString(rd, "", flatten.DotStyle)
	//if err != nil {
	//    fmt.Println("Flattening error")
	//}
	//olog.Print(flat)
	//tableprinter.Print(os.Stdout, rd)

	//printer := tableprinter.New(os.Stdout)
	//printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	//printer.CenterSeparator = "│"
	//printer.ColumnSeparator = "│"
	//printer.RowSeparator = "─"
	//printer.HeaderBgColor = tablewriter.BgBlackColor // set header background color for all headers.
	//printer.HeaderFgColor = tablewriter.FgGreenColor // set header foreground color for all headers.

	var phost []PstoreHost
	if hgroup != "" {
		//printer.Print(rd)
		olog.Print(rd)
		for _, r := range rd {
			for _, h := range r.Hosts {
				phost = append(phost, h)
			}
		}
		//printer.Print(phost)
		olog.Print(phost)
	} else {
		//printer.Print(rd)
		olog.Print(rd)
	}

}
*/
