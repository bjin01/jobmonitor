package spmigration

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func Get_salt_online_Minions_in_Group(sessionkey *auth.SumaSessionKey, minion_list []string, groupsdata *Migration_Groups) []string {
	var offline_minion_list []string
	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Target_List = minion_list
	saltdata.Login()
	saltdata.Run_Manage_Status()
	if len(saltdata.Offline_Minions) > 0 {
		offline_minion_list = saltdata.Offline_Minions
		fmt.Printf("Salt offline minions: %v\n", offline_minion_list)
	} else {
		offline_minion_list = []string{}
	}

	if len(saltdata.Online_Minions) > 0 {
		fmt.Printf("Salt online minions: %v\n", saltdata.Online_Minions)
	} else {
		fmt.Printf("Salt online minions is empty\n")
	}
	return offline_minion_list
}

func string_array_contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
