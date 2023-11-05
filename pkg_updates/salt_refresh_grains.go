package pkg_updates

import (
	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
	"gorm.io/gorm"
)

func Salt_Refresh_Grains_New(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB) {

	all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}
	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword

	for _, minion := range all_minions {
		if minion.Minion_Status == "Online" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Minion_Remarks", "")
		}
	}

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		refresh_grains_return := saltdata.Saltutil_Refresh_Grains()
		if len(refresh_grains_return) > 0 {
			logger.Infof("Minions saltutil.refresh_grains return: %v\n", refresh_grains_return)
		}
	}

}
