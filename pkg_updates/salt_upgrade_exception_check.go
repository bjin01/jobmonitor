package pkg_updates

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
	"gorm.io/gorm"
)

func Salt_No_Upgrade_Exception_Check_New(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB) {
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
	saltdata.Salt_already_patched_exception_key = groupsdata.Salt_already_patched_exception_key
	saltdata.Salt_already_patched_exception_value = groupsdata.Salt_already_patched_exception_value

	all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}
	for _, minion := range all_minions {
		//logger.Infof("Minion %s is ready for disk space check\n", minion.Minion_Name)
		if minion.Minion_Status == "Online" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
		}
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		//saltdata.Run_Refresh_Grains()
		disqualified_minions := saltdata.Run_No_Upgrade_Grains_Check()
		logger.Infof("Salt_No_Upgrade_Exception_Check: %v\n", disqualified_minions)

		exception_keys := []struct {
			Key     string
			Value   string
			Subject string
		}{
			{
				Key:     groupsdata.Salt_no_upgrade_exception_key,
				Value:   groupsdata.Salt_no_upgrade_exception_value,
				Subject: "No_upgrade exception",
			},
			{
				Key:     groupsdata.Salt_already_patched_exception_key,
				Value:   groupsdata.Salt_already_patched_exception_value,
				Subject: "Already_patched exception",
			},
		}

		for _, ex := range exception_keys {
			if ex.Key == "" || ex.Value == "" {
				continue
			}
			// Set the current exception key/value in saltdata
			saltdata.Salt_no_upgrade_exception_key = ex.Key
			saltdata.Salt_no_upgrade_exception_value = ex.Value

			disqualified_minions := saltdata.Run_No_Upgrade_Grains_Check()
			for _, minion := range all_minions {
				if string_array_contains(disqualified_minions, minion.Minion_Name) {
					logger.Infof("Minion %s has %s and is disqualified\n", minion.Minion_Name, ex.Subject)
					subject := ex.Subject
					body := fmt.Sprintf("%s for minion found: %s", ex.Subject, minion.Minion_Name)
					Add_Note(sessionkey, minion.Minion_ID, subject, body)
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Minion_Remarks", ex.Subject+" is true")
				}
			}
		}

	}

}

func Salt_No_Upgrade_Exception_Check_New_by_List(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, minion_list []Minion_Data, db *gorm.DB) {
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

	for _, minion := range minion_list {
		//logger.Infof("Minion %s is ready for disk space check\n", minion.Minion_Name)
		if minion.Minion_Status == "Online" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
		}
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		disqualified_minions := saltdata.Run_No_Upgrade_Grains_Check()
		logger.Debugf("Salt_No_Upgrade_Exception_Check: %v\n", disqualified_minions)

		for _, minion := range minion_list {
			if string_array_contains(disqualified_minions, minion.Minion_Name) {
				logger.Infof("Minion %s has no_upgrade exception and is disqualified\n", minion.Minion_Name)
				subject := "No_upgrade exception"
				body := fmt.Sprintf("No_upgrade exception for minion found: %s", minion.Minion_Name)
				Add_Note(sessionkey, minion.Minion_ID, subject, body)
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Minion_Remarks", "No_Upgrade_Exception is true")
			}
		}
	}
}
