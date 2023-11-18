package pkg_updates

import (
	"context"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"gorm.io/gorm"
)

func Start_Workflow(ctx context.Context, sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB, health *bool, deadline *time.Time) {

	gr := getGoroutineID()
	logger.Infof("Start_Workflow: %s: Goroutine ID %d\n", groupsdata.T7User, gr)
	sleep_between_steps := 2 * time.Second
	for time.Now().Before(*deadline) {
		if *health == false {
			logger.Infof("Check_Jobs can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Debugf("Package Upgrade Workflow err: %s\n", err)
			}
			logger.Infof("Package Upgrade Workflow context: finished %s\n", groupsdata.Ctx_ID)
			return
		default:
			logger.Infof("Package Upgrade Workflow: running %s\n", groupsdata.Ctx_ID)
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
					//logger.Infof("Minion %s is ready for workflow\n", minion.Minion_Name)
					*qualified_minions = append(*qualified_minions, minion)
				}
			}
		}

		if len(wf) != 0 {
			for _, w := range wf {

				if w.Name == "assign_channels" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Assign_Channels(sessionkey, groupsdata, db, wf, *qualified_minions, "assign_channels")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_updates" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Update_packages(sessionkey, db, wf, *qualified_minions, "package_updates")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_update_reboot" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Reboot(sessionkey, db, wf, *qualified_minions, "package_update_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_refresh" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Refresh_Packages(sessionkey, db, wf, *qualified_minions, "package_refresh")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "waiting" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Waiting_Stage(db, wf, *qualified_minions, "waiting")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_dryrun" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_dryrun")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_dryrun", true)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_run" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_run")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_run", false)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_reboot" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Reboot(sessionkey, db, wf, *qualified_minions, "spmigration_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_package_refresh" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Refresh_Packages(sessionkey, db, wf, *qualified_minions, "spmigration_package_refresh")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "post_migration" {
					logger.Debugf("Workflow %s - %s - %s\n", w.Name, groupsdata.T7User, groupsdata.Groups[0])
					Post_Migration(db, groupsdata, wf, *qualified_minions, "post_migration")
				}
			}
		}
		logger.Infof("Start_Workflow: continue after 20 seconds\n")
		time.Sleep(20 * time.Second)
	}
	logger.Infof("Workflow final deadline reached. Exiting.\n")
	return
}
