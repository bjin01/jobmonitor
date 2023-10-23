package spmigration

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func (m *Target_Minions) Salt_No_Upgrade_Exception_Check(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {
	logger.Infof("Salt_No_Upgrade_Exception_Check\n")
	//logger.Infof("groupsdata: %v: %v\n", groupsdata.Salt_no_upgrade_exception_key, groupsdata.Salt_no_upgrade_exception_value)
	if groupsdata.Salt_no_upgrade_exception_key == "" || groupsdata.Salt_no_upgrade_exception_value == "" {
		logger.Infof("Salt no_upgrade grains key and or value is not provided. Skipping.\n")
		return
	}

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Salt_no_upgrade_exception_key = groupsdata.Salt_no_upgrade_exception_key
	saltdata.Salt_no_upgrade_exception_value = groupsdata.Salt_no_upgrade_exception_value

	for _, minion := range m.Minion_List {
		//logger.Infof("Minion %s is ready for disk space check\n", minion.Minion_Name)
		saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	for _, minion := range m.No_Targets_Minions {
		saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		//saltdata.Run_Refresh_Grains()
		disqualified_minions := saltdata.Run_No_Upgrade_Grains_Check()
		if len(disqualified_minions) > 0 {
			m.No_Upgrade_Exceptions = disqualified_minions
			//logger.Infof("Minions disqualified by no_upgrade exception check: %v\n", disqualified_minions)
			newMinionList := new([]Minion_Data)
			for _, minion := range m.Minion_List {
				if !string_array_contains(disqualified_minions, minion.Minion_Name) {
					*newMinionList = append(*newMinionList, minion)
				} else {
					logger.Infof("Minion %s has no_upgrade exception and is isqualified\n", minion.Minion_Name)
					subject := "No_upgrade exception"
					body := fmt.Sprintf("No_upgrade exception for minion found: %s %s", minion.Minion_Name, m.Suma_Group)
					Add_Note(sessionkey, minion.Minion_ID, subject, body)
				}
			}
			if len(*newMinionList) > 0 {
				logger.Infof("Minion list after no_upgrade exception check: %v\n", newMinionList)
				m.Minion_List = *newMinionList

			}

			if len(*newMinionList) == 0 {
				logger.Infof("All minions have been disqualified by no_upgrade exception check. Exiting.\n")
				m.Minion_List = []Minion_Data{}
				return
			}

		} else {
			logger.Infof("All minions passed no_upgrade exception check.\n")
		}
	}

}
