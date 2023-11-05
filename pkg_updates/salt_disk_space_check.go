package pkg_updates

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
	"gorm.io/gorm"
)

func Salt_Disk_Space_Check_New(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB) {

	if groupsdata.Salt_diskspace_grains_key == "" || groupsdata.Salt_diskspace_grains_value == "" {
		logger.Infof("Salt disk space check is not configured. Skipping.\n")
		return
	}

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.Salt_diskspace_grains_key = groupsdata.Salt_diskspace_grains_key
	saltdata.Salt_diskspace_grains_value = groupsdata.Salt_diskspace_grains_value

	all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	for _, minion := range all_minions {
		if minion.Minion_Status == "Online" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
		}
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		//saltdata.Run_Refresh_Grains()
		disqualified_minions := saltdata.Run_Disk_Space_Check()

		logger.Infof("Minions disqualified by disk space check: %v\n", disqualified_minions)
		for _, minion := range all_minions {
			if string_array_contains(disqualified_minions, minion.Minion_Name) {
				logger.Infof("Minion %s is disk space check disqualified\n", minion.Minion_Name)
				subject := "btrfs disqualified"
				note := fmt.Sprintf("/ has less than 2GB free space. %s", minion.Minion_Name)
				Add_Note(sessionkey, minion.Minion_ID, subject, note)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Minion_Remarks", "btrfs disk space check disqualified. <2GB")
			} /* else {
				logger.Infof("Minion %s passed disk space check\n", minion.Minion_Name)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Minion_Remarks", "")
			} */
		}

	}

}
