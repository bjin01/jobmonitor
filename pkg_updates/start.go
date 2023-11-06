package pkg_updates

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"gorm.io/gorm"
)

func Start_Workflow(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB, health *bool) {
	// TODO - Add your code here
	deadline := time.Now().Add(time.Duration(60) * time.Minute)
	sleep_between_steps := 2 * time.Second
	for time.Now().Before(deadline) {
		if *health == false {
			logger.Infof("Check_Jobs can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		wf, err := Get_Workflow_From_DB(db)
		if err != nil {
			logger.Errorf("failed to get workflow from database")
			return
		}

		all_minions, err := GetAll_Minions_From_DB(db)
		if err != nil {
			logger.Errorf("failed to connect database")
			return
		}

		qualified_minions := new([]Minion_Data)
		if len(all_minions) > 0 {
			for _, minion := range all_minions {
				if minion.Minion_Status == "Online" && minion.Minion_Remarks == "" {
					logger.Infof("Minion %s is ready for workflow\n", minion.Minion_Name)
					*qualified_minions = append(*qualified_minions, minion)
				}
			}
		}

		if len(wf) != 0 {
			for _, w := range wf {

				if w.Name == "assigne_channels" {
					logger.Infof("Workflow %s\n", w.Name)
					Assign_Channels(sessionkey, groupsdata, db, wf, *qualified_minions, "assigne_channels")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_updates" {
					logger.Infof("Workflow %s\n", w.Name)
					Update_packages(sessionkey, db, wf, *qualified_minions, "package_updates")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_update_reboot" {
					logger.Infof("Workflow %s\n", w.Name)
					Reboot(sessionkey, db, wf, *qualified_minions, "package_update_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_refresh" {
					logger.Infof("Workflow %s\n", w.Name)
					Refresh_Packages(sessionkey, db, wf, *qualified_minions, "package_refresh")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "waiting" {
					logger.Infof("Workflow %s\n", w.Name)
					Waiting_Stage(db, wf, *qualified_minions, "waiting")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_dryrun" {
					logger.Infof("Workflow %s\n", w.Name)
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_dryrun")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_dryrun", true)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_run" {
					logger.Infof("Workflow %s\n", w.Name)
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_dryrun")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_run", false)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_reboot" {
					logger.Infof("Workflow %s\n", w.Name)
					Reboot(sessionkey, db, wf, *qualified_minions, "spmigration_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "post_migration" {
					logger.Infof("Workflow %s\n", w.Name)
					Post_Migration(db, wf, *qualified_minions, "post_migration")
				}
			}
		}
		logger.Infof("Start_Workflow: continue after 20 seconds\n")
		time.Sleep(20 * time.Second)
	}
}
