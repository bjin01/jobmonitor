package spmigration

import "github.com/bjin01/jobmonitor/auth"

func (m *Target_Minions) Detect_Online_Minions(sessionkey *auth.SumaSessionKey, list []Minion_Data, groupsdata *Migration_Groups) (online_minions []Minion_Data, offline_minions []string) {

	unique_minions := []Minion_Data{}
	var Online_Minion_List []Minion_Data
	var offline_minion_list []string

	for _, minion := range list {
		exists := false
		// Check if the Minion_ID is already in the map
		for _, minion2 := range m.Minion_List {
			if minion.Minion_ID == minion2.Minion_ID {
				// If it is, don't add it again
				logger.Infof("Minion %s already exists in Minion_List\n", minion.Minion_Name)
				exists = true
			}
		}
		if !exists {
			//logger.Infof("Adding Minion to Minion_List: %s\n", minion.Minion_Name)
			unique_minions = append(unique_minions, minion)
		}

	}

	if len(unique_minions) > 0 {
		var salt_minion_list []string
		for _, minion := range unique_minions {
			salt_minion_list = append(salt_minion_list, minion.Minion_Name)
		}
		offline_minions = Get_salt_online_Minions_in_Group(sessionkey, salt_minion_list, groupsdata)

		for _, minion2 := range unique_minions {
			if !string_array_contains(offline_minions, minion2.Minion_Name) {
				Online_Minion_List = append(Online_Minion_List, minion2)
			} else {
				logger.Infof("Minion %s is offline\n", minion2.Minion_Name)
				subject := "minion is offline"
				body := "minion is offline"
				Add_Note(sessionkey, minion2.Minion_ID, subject, body)
				offline_minion_list = append(offline_minion_list, minion2.Minion_Name)
			}
		}

	}
	return Online_Minion_List, offline_minion_list
}
