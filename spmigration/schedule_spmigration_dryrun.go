package spmigration

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ScheduleSPMigrationDryRun_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Sid        []int  `xmlrpc:"sid"`
}

type ListMigrationTarget_Request struct {
	Sessionkey                          string `xmlrpc:"sessionKey"`
	Sid                                 int    `xmlrpc:"sid"`
	ExcludeTargetWhereMissingSuccessors bool   `xmlrpc:"excludeTargetWhereMissingSuccessors"`
}

type ListMigrationTarget_Response struct {
	Result []struct {
		Ident    string `xmlrpc:"ident,omitempty"`
		Friendly string `xmlrpc:"friendly,omitempty"`
	}
}

type ListAllChannels_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
}

type ListAllChannels_Response struct {
	Result []struct {
		Id            int    `xmlrpc:"id,omitempty"`
		Name          string `xmlrpc:"name,omitempty"`
		Label         string `xmlrpc:"label,omitempty"`
		Arch_name     string `xmlrpc:"arch_name,omitempty"`
		Provider_name string `xmlrpc:"provider_name,omitempty"`
		Packages      int    `xmlrpc:"packages,omitempty"`
		Systems       int    `xmlrpc:"systems,omitempty"`
	}
}

func (t *Target_Minions) ListMigrationTarget(sessionkey *auth.SumaSessionKey, UserData *Migration_Groups) {
	method := "system.listMigrationTargets"
	allchannels := List_All_Channels(sessionkey)
	for i, minion := range t.Minion_List {
		t.Minion_List[i].Target_base_channel = UserData.Target_base_channel
		var params ListMigrationTarget_Request
		params.Sessionkey = sessionkey.Sessionkey
		params.Sid = minion.Minion_ID
		params.ExcludeTargetWhereMissingSuccessors = true
		buf, err := gorillaxml.EncodeClientRequest(method, &params)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		//fmt.Printf("buffer: %s\n", string(buf))
		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		/* responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}
		fmt.Printf("responseBody: %s\n", responseBody) */
		reply := new(ListMigrationTarget_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			log.Fatalf("Decode ListMigrationTarget response body failed: %s\n", err)
		}
		for _, target := range reply.Result {
			fmt.Printf("%s:\n", minion.Minion_Name)
			fmt.Printf("Ident: %s\n", target.Ident)
			returned_map := Convert_String_to_maps(target.Friendly)
			//fmt.Printf("Target base channel: %s\n", UserData.Target_base_channel)
			for _, channel := range allchannels.Result {
				fmt.Printf("Channel name: %s\n", channel.Name)
				for key, value := range returned_map {
					if key == "base" {
						if value == channel.Name {
							//t.Minion_List[i].Target_base_channel_id = channel.Id
							fmt.Printf("Found target base channel name: %s\n", channel.Name)
						}
					}
				}
			}
		}
		//log.Printf("ListMigrationTarget: %s\n", reply.Result)

	}
}

func Convert_String_to_maps(mystring string) map[string]string {
	//mystring := "[base: SUSE Linux Enterprise Server for SAP Applications 15 SP5 x86_64, addon: Desktop Applications Module 15 SP5 x86_64, SUSE Linux Enterprise Live Patching 15 SP5 x86_64, Web and Scripting Module 15 SP5 x86_64, Basesystem Module 15 SP5 x86_64, SAP Applications Module 15 SP5 x86_64, Server Applications Module 15 SP5 x86_64, SUSE Manager Client Tools for SLE 15 x86_64, Python 3 Module 15 SP5 x86_64, SUSE Linux Enterprise High Availability Extension 15 SP5 x86_64]"

	// Remove brackets from the string
	mystring = strings.TrimPrefix(mystring, "[")
	mystring = strings.TrimSuffix(mystring, "]")

	// Split the string into key-value pairs
	pairs := strings.Split(mystring, ", ")

	// Create a map to store the key-value pairs
	m := make(map[string]string)

	// Iterate over each pair and populate the map
	for _, pair := range pairs {
		kv := strings.Split(pair, ": ")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[key] = value
		}
	}

	// Print the resulting map
	for key, value := range m {
		fmt.Printf("%s: %s\n", key, value)
	}
	return m
}

func Convert_String_IntSlices(mystring string) []int {

	// Remove brackets from the string
	mystring = strings.TrimPrefix(mystring, "[")
	mystring = strings.TrimSuffix(mystring, "]")

	// Split the string into individual integers
	intStrs := strings.Split(mystring, ",")

	// Create a slice to store the integers
	intSlice := make([]int, 0, len(intStrs))

	// Convert each string to an integer and append it to the slice
	for _, s := range intStrs {
		if i, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			intSlice = append(intSlice, i)
		}
	}

	// Print the resulting slice
	//fmt.Println(intSlice)
	return intSlice
}

func List_All_Channels(sessionkey *auth.SumaSessionKey) *ListAllChannels_Response {
	method := "channel.listAllChannels"
	var params ListAllChannels_Request
	params.Sessionkey = sessionkey.Sessionkey

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", string(buf))
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */
	reply := new(ListAllChannels_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Fatalf("Decode List_All_Channels response body failed: %s\n", err)
	}
	//fmt.Printf("List_All_Vendor_Channels: %v\n", reply.Result)
	return reply
}
