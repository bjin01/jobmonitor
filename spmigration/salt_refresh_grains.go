package spmigration

import (
	"log"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func (m *Target_Minions) Salt_Refresh_Grains(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword

	for _, minion := range m.Minion_List {
		saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		refresh_grains_return := saltdata.Saltutil_Refresh_Grains()
		if len(refresh_grains_return) > 0 {
			log.Printf("Minions saltutil.refresh_grains return: %v\n", refresh_grains_return)
		}
	}

}
