package pkg_updates

import "github.com/bjin01/jobmonitor/auth"

func Detect_Online_Minions(sessionkey *auth.SumaSessionKey, list []Minion_Data, groupsdata *Update_Groups) []Minion_Data {

	unique_minions := []Minion_Data{}
	var Return_Minion_List []Minion_Data

	for _, minion := range list {
		exists := false
		// Check if the Minion_ID is already in the map
		for _, unique_minion := range unique_minions {
			if minion.Minion_ID == unique_minion.Minion_ID {
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
		offline_minions := Get_salt_online_Minions_in_Group_New(sessionkey, salt_minion_list, groupsdata)

		for _, minion2 := range unique_minions {
			if !string_array_contains(offline_minions, minion2.Minion_Name) {
				minion2.Minion_Status = "Online"
				Return_Minion_List = append(Return_Minion_List, minion2)
			} else {
				logger.Infof("Minion %s is offline\n", minion2.Minion_Name)
				subject := "minion is offline"
				body := "minion is offline"
				Add_Note(sessionkey, minion2.Minion_ID, subject, body)
				minion2.Minion_Status = "Offline"
				Return_Minion_List = append(Return_Minion_List, minion2)
			}
		}

	}
	return Return_Minion_List
}
