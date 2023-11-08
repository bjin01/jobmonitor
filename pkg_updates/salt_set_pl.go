package pkg_updates

import (
	"github.com/bjin01/jobmonitor/saltapi"
	"gorm.io/gorm"
)

func Salt_Set_Patch_Level(groupsdata *Update_Groups, minion_list []string, stage string, db *gorm.DB) {

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

	saltdata.Online_Minions = minion_list

	saltdata_refresh_grains := new(saltapi.Salt_Data)
	saltdata_refresh_grains.SaltMaster = groupsdata.SaltMaster_Address
	saltdata_refresh_grains.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata_refresh_grains.Username = groupsdata.SaltUser
	saltdata_refresh_grains.Password = groupsdata.SaltPassword

	saltdata_refresh_grains.Online_Minions = minion_list

	if len(saltdata.Online_Minions) > 0 {
		saltdata.Login()
		set_pl_return := saltdata.Run_Set_Patch_Level()
		if len(set_pl_return) > 0 {
			logger.Infof("Minions set patch level done: %d returned\n", len(set_pl_return))
		}

		//must include this to refresh grains otherwise the patch level will not be updated

		saltdata_refresh_grains.Login()
		refresh_grains_return := saltdata_refresh_grains.Saltutil_Refresh_Grains()
		if len(refresh_grains_return) > 0 {
			logger.Infof("Minions saltutil.refresh_grains return: %v\n", refresh_grains_return)
		}
	}

}
