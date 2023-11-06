package pkg_updates

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAll_Minions_From_DB(db *gorm.DB) ([]Minion_Data, error) {
	var minion_data []Minion_Data
	err := db.Preload(clause.Associations).Find(&minion_data).Error
	//err := db.Model(&grp).Preload("Posts").Find(&grp).Error
	return minion_data, err
}

func Get_Workflow_From_DB(db *gorm.DB) ([]Workflow_Step, error) {
	var workflow_step []Workflow_Step
	err := db.Preload(clause.Associations).Find(&workflow_step).Error
	//err := db.Model(&grp).Preload("Posts").Find(&grp).Error
	return workflow_step, err
}
