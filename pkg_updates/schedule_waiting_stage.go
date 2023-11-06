package pkg_updates

import (
	"gorm.io/gorm"
)

func Waiting_Stage(db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}

		//fmt.Printf("-----------Query DB waiting stage %d\n", result.RowsAffected)
		logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", 3)
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
			logger.Infof("Minion %s starts %s stage.\n", minion.Minion_Name, stage)
		}
	}

}
