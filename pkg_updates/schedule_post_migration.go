package pkg_updates

import (
	"gorm.io/gorm"
)

func Post_Migration(db *gorm.DB, groupsdata *Update_Groups, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	post_minion_list := []string{}
	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}

		//fmt.Printf("-----------Query DB waiting stage %d\n", result.RowsAffected)
		//logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", 3)
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
			logger.Debugf("Minion %s starts %s stage.\n", minion.Minion_Name, stage)
			post_minion_list = append(post_minion_list, minion.Minion_Name)

		}
	}
	go Salt_Run_Post_State(groupsdata, post_minion_list)
	go Salt_Set_Patch_Level(groupsdata, wf, post_minion_list, stage, db)

}
