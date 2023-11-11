package main

import (
	"github.com/bjin01/jobmonitor/pkg_updates"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Pkg_update_get_minion_from_db(filename string, minion_name string) (pkg_updates.Minion_Data, error) {
	var minion pkg_updates.Minion_Data
	logger.Infof("Use sqlite database: %s and query single minion %s\n", filename, minion_name)
	db, err := gorm.Open(gorm.Dialector(&sqlite.Dialector{DSN: filename}),
		&gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		logger.Errorf("failed to connect database")
		return minion, err
	}

	minion.Minion_Name = minion_name
	db.Preload("Target_Optional_Channels").Preload("Minion_Groups").Where("Minion_Name = ?", minion_name).Find(&minion)

	//Below code can be enabled to list all optional channels
	/* var optchannels []pkg_updates.OptionalChannels
	result := db.Find(&optchannels)
	if result.RowsAffected > 0 {
		logger.Infof("Found %d optional channels\n", result.RowsAffected)
		for _, oc := range optchannels {
			logger.Debugf("Optional channel: %s ID: %d\n", oc.Channel_Label, oc.ID)
		}
	} */

	return minion, nil
}
