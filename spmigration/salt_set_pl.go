package spmigration

import (
	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func (m *Target_Minions) Salt_Set_Patch_Level(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {

	if groupsdata.Patch_Level == "" {
		logger.Infof("Patch Level is not provided. Skipping.\n")
		return
	}

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Patch_Level = groupsdata.Patch_Level

	for _, minion := range m.Minion_List {
		if minion.Migration_Stage == "Product Migration" && minion.Migration_Stage_Status == "Completed" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
		}
		//saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	for _, minion := range m.No_Targets_Minions {
		saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		set_pl_return := saltdata.Run_Set_Patch_Level()
		if len(set_pl_return) > 0 {
			logger.Infof("Minions set patch level done: %d returned\n", len(set_pl_return))
		}
	}

}
