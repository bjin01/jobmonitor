package pkg_updates

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func Get_salt_online_Minions_in_Group_New(sessionkey *auth.SumaSessionKey, minion_list []string, groupsdata *Update_Groups) []string {
	var offline_minion_list []string
	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Target_List = minion_list

	timeout := fmt.Sprintf("timeout=%d", groupsdata.Timeout)
	gather_job_timeout := fmt.Sprintf("gather_job_timeout=%d", groupsdata.GatherJobTimeout)
	logger.Debugf("salt-run manage.status timeout: %s\n", timeout)
	logger.Debugf("salt-run manage.status gather_job_timeout: %s\n", gather_job_timeout)
	saltdata.Arg = append(saltdata.Arg, timeout)
	saltdata.Arg = append(saltdata.Arg, gather_job_timeout)

	saltdata.Login()
	saltdata.Run_Manage_Status()
	if len(saltdata.Offline_Minions) > 0 {
		offline_minion_list = saltdata.Offline_Minions
		logger.Infof("Salt offline minions: %v\n", offline_minion_list)
	} else {
		offline_minion_list = []string{}
	}

	if len(saltdata.Online_Minions) > 0 {
		logger.Infof("Salt online minions: %v\n", saltdata.Online_Minions)
	} else {
		logger.Infof("Salt online minions is empty\n")
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
