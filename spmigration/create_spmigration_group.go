package spmigration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Create_SPMigration_Group_Request struct {
	Sessionkey  string `xmlrpc:"sessionKey"`
	Name        string `xmlrpc:"name"`
	Description string `xmlrpc:"description"`
}

type Create_SPMigration_Group_Response struct {
	Server_group struct {
		Id           int    `xmlrpc:"id"`
		Name         string `xmlrpc:"name"`
		Description  string `xmlrpc:"description"`
		Org_id       int    `xmlrpc:"org_id"`
		System_count int    `xmlrpc:"system_count"`
	} `xmlrpc:"server group"`
}

type AddOrRemoveSystems_Request struct {
	Sessionkey      string `xmlrpc:"sessionKey"`
	SystemGroupName string `xmlrpc:"systemGroupName"`
	ServerIds       []int  `xmlrpc:"serverIds"`
	Add             bool   `xmlrpc:"add"`
}

type AddOrRemoveSystems_Response struct {
	Result_ID int
}

func parseFilename(input string) string {
	parts := strings.Split(input, "/")
	filename := parts[len(parts)-1]

	if strings.HasPrefix(filename, "spmigration_") && strings.HasSuffix(filename, ".yaml") {
		return filename[0 : len(filename)-5]
	}

	return ""
}

func (t *Target_Minions) SPMigration_Group(sessionkey *auth.SumaSessionKey, UserData *Migration_Groups) {
	if UserData.T7User == "" {
		t.Suma_Group = fmt.Sprintf("SPMigration_%s", time.Now().Format("20060102150405"))
	} else {
		t.Suma_Group = fmt.Sprintf("SPMigration_%s_%s", UserData.T7User, time.Now().Format("20060102150405"))
	}

	group_name := ""
	if t.Tracking_file_name != "" {
		group_name = parseFilename(t.Tracking_file_name)
	}

	if group_name != "" {
		t.Suma_Group = group_name
	} else {
		if UserData.T7User == "" {
			t.Suma_Group = fmt.Sprintf("SPMigration_%s", time.Now().Format("20060102150405"))
		} else {
			t.Suma_Group = fmt.Sprintf("SPMigration_%s_%s", UserData.T7User, time.Now().Format("20060102150405"))
		}
	}

	if len(t.Minion_List) == 0 {
		log.Println("Minion list is empty, no need to create SPMigration group")
		return
	}

	t.Create_SPMigration_Group(sessionkey, UserData)
	t.Add_Systems_To_SPMigration_Group(sessionkey)

}

func (t *Target_Minions) Add_Systems_To_SPMigration_Group(sessionkey *auth.SumaSessionKey) {
	method := "systemgroup.addOrRemoveSystems"
	var serverids []int
	var serveridsMap = make(map[int]bool)
	for _, minion := range t.Minion_List {
		if _, exists := serveridsMap[minion.Minion_ID]; !exists {
			// If it's not in the map, add it to the slice and mark it as seen in the map
			serverids = append(serverids, minion.Minion_ID)
			serveridsMap[minion.Minion_ID] = true
		}
	}
	params := AddOrRemoveSystems_Request{
		Sessionkey:      sessionkey.Sessionkey,
		SystemGroupName: t.Suma_Group,
		ServerIds:       serverids,
		Add:             true,
	}
	buf, err := gorillaxml.EncodeClientRequest(method, &params)
	if err != nil {
		log.Fatalf("Encoding systemgroup.addOrRemoveSystems error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("response: %s\n", responseBody) */

	reply := new(AddOrRemoveSystems_Response)

	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Printf("Decode AddOrRemoveSystems_Response response body failed: %s\n", err)
	}
	if reply.Result_ID == 1 {
		log.Printf("%d Minions added to group %s\n", len(serverids), t.Suma_Group)
	} else {
		log.Printf("Failed to add minions to group %s\n", t.Suma_Group)
	}
}

func (t *Target_Minions) Create_SPMigration_Group(sessionkey *auth.SumaSessionKey,
	UserData *Migration_Groups) {
	method := "systemgroup.create"
	params := Create_SPMigration_Group_Request{
		Sessionkey:  sessionkey.Sessionkey,
		Name:        t.Suma_Group,
		Description: t.Suma_Group,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
	if err != nil {
		log.Fatalf("Encoding systemgroup.create error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("response: %v\n", resp.Body)

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("response: %s\n", responseBody) */
	reply := new(Create_SPMigration_Group_Response)

	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Fatalf("Decode Create_SPMigration_Group_Response response body failed: %s\n", err)
	}

	log.Printf("SPMigration group %s created\n", reply.Server_group.Name)
}
