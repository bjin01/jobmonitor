package pkg_updates

import (
	"context"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Start_Workflow(ctx context.Context, sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB, health *bool, deadline *time.Time) {

	gr := getGoroutineID()
	logger.WithFields(logrus.Fields{
		"goroutine": gr,
		"username":  groupsdata.T7User,
	}).Infoln("Workflow started")

	sleep_between_steps := 2 * time.Second
	for time.Now().Before(*deadline) {
		if !*health {
			logger.WithFields(logrus.Fields{
				"goroutine": gr,
				"username":  groupsdata.T7User,
			}).Infof("Check_Jobs can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.")
			time.Sleep(125 * time.Second)
			continue
		}

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.WithFields(logrus.Fields{
					"goroutine": gr,
					"username":  groupsdata.T7User,
				}).Debugf("Package Upgrade Workflow err: %s", err)
			}
			logger.WithFields(logrus.Fields{
				"goroutine": gr,
				"username":  groupsdata.T7User,
			}).Infof("Package Upgrade Workflow context: finished %s", groupsdata.Ctx_ID)
			return
		default:
			logger.WithFields(logrus.Fields{
				"goroutine": gr,
				"username":  groupsdata.T7User,
			}).Infof("Package Upgrade Workflow: running %s", groupsdata.Ctx_ID)
		}

		wf, err := Get_Workflow_From_DB(db)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"goroutine": gr,
				"username":  groupsdata.T7User,
			}).Errorf("failed to get workflow from database")
			return
		}

		all_minions, err := GetAll_Minions_From_DB(db)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"goroutine": gr,
				"username":  groupsdata.T7User,
			}).Errorf("failed to connect database")
			return
		}

		qualified_minions := new([]Minion_Data)
		if len(all_minions) > 0 {
			for _, minion := range all_minions {
				if minion.Minion_Status == "Online" && minion.Minion_Remarks == "" {
					//logger.Infof("Minion %s is ready for workflow", minion.Minion_Name)
					*qualified_minions = append(*qualified_minions, minion)
				}
			}
		}

		if len(wf) != 0 {
			for _, w := range wf {

				if w.Name == "assign_channels" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Assign_Channels(sessionkey, groupsdata, db, wf, *qualified_minions, "assign_channels")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_updates" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Update_packages(sessionkey, db, wf, *qualified_minions, "package_updates")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_update_reboot" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Reboot(sessionkey, db, wf, *qualified_minions, "package_update_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "package_refresh" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Refresh_Packages(sessionkey, db, wf, *qualified_minions, "package_refresh")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "waiting" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Waiting_Stage(db, wf, *qualified_minions, "waiting")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_dryrun" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_dryrun")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_dryrun", true)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_run" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					ListMigrationTarget(sessionkey, groupsdata, db, wf, *qualified_minions, "spmigration_run")
					SPMigration(sessionkey, db, wf, *qualified_minions, "spmigration_run", false)
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_reboot" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Reboot(sessionkey, db, wf, *qualified_minions, "spmigration_reboot")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "spmigration_package_refresh" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Refresh_Packages(sessionkey, db, wf, *qualified_minions, "spmigration_package_refresh")
				}
				time.Sleep(sleep_between_steps)

				if w.Name == "post_migration" {
					logger.WithFields(logrus.Fields{
						"goroutine": gr,
						"username":  groupsdata.T7User,
					}).Debugf("Workflow %s - %s", w.Name, groupsdata.Groups[0])
					Post_Migration(db, groupsdata, wf, *qualified_minions, "post_migration")
				}
			}
		}
		logger.WithFields(logrus.Fields{
			"goroutine": gr,
			"username":  groupsdata.T7User,
		}).Infof("Start_Workflow: continue after 20 seconds")
		time.Sleep(20 * time.Second)
	}
	logger.WithFields(logrus.Fields{
		"goroutine": gr,
		"username":  groupsdata.T7User,
	}).Infof("Workflow final deadline reached. Exiting.")
}
