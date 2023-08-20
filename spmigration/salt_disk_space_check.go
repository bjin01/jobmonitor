package spmigration

import (
	"log"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func (m *Target_Minions) Salt_Disk_Space_Check(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {

	if groupsdata.Salt_diskspace_grains_key == "" || groupsdata.Salt_diskspace_grains_value == "" {
		log.Printf("Salt disk space check is not configured. Skipping.\n")
		return
	}

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Salt_diskspace_grains_key = groupsdata.Salt_diskspace_grains_key
	saltdata.Salt_diskspace_grains_value = groupsdata.Salt_diskspace_grains_value

	for _, minion := range m.Minion_List {
		//fmt.Printf("Minion %s is ready for disk space check\n", minion.Minion_Name)
		saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		//saltdata.Run_Refresh_Grains()
		disqualified_minions := saltdata.Run_Disk_Space_Check()
		if len(disqualified_minions) > 0 {
			m.Disk_Check_Disqualified = disqualified_minions
			log.Printf("Minions disqualified by disk space check: %v\n", disqualified_minions)
			newMinionList := new([]Minion_Data)
			for _, minion := range m.Minion_List {
				if !string_array_contains(disqualified_minions, minion.Minion_Name) {
					*newMinionList = append(*newMinionList, minion)
				} else {
					log.Printf("Minion %s is disk space check disqualified\n", minion.Minion_Name)
				}
			}
			if len(*newMinionList) > 0 {
				log.Printf("Minion list after disk space check: %v\n", newMinionList)
				m.Minion_List = *newMinionList

			}

			if len(*newMinionList) == 0 {
				log.Printf("All minions have been disqualified by disk space check. Exiting.\n")
				m.Minion_List = []Minion_Data{}
				return
			}

		} else {
			log.Printf("All minions passed disk space check.\n")
		}
	}

}