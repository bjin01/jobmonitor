package pkg_updates

import (
	"context"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
	"gorm.io/gorm"
)

type Job_Data struct {
	Scheduler         string
	Name              string
	CompletedSystems  int
	FailedSystems     int
	InProgressSystems int
	Id                int
	Type              string
	Earliest          time.Time
}

func Check_Jobs(ctx context.Context, groupsdata *Update_Groups, sessionkey *auth.SumaSessionKey, health *bool, db *gorm.DB, deadline *time.Time) {

	//deadline := time.Now().Add(time.Duration(60) * time.Minute)
	gr := getGoroutineID()
	logger.Infof("Check_Jobs: Goroutine ID %d", gr)

	for time.Now().Before(*deadline) {
		if *health == false {
			logger.Infof("Check_Jobs can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.")
			time.Sleep(125 * time.Second)
			continue
		}

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Debugf("Check_Jobs err: %s %s", err, groupsdata.Ctx_ID)
			}
			logger.Infof("Check_Jobs: finished %s", groupsdata.Ctx_ID)
			return
		default:
			logger.Infof("Check_Jobs: running %s", groupsdata.Ctx_ID)
		}

		all_minions, err := GetAll_Minions_From_DB(db)
		if err != nil {
			logger.Errorf("failed to connect database")
			return
		}

		var joblist schedules.ListJobs
		joblist.Found_Pending_Jobs = false

		/* if len(all_minions) > 0 {
			joblist.GetPendingjobs(sessionkey)
			joblist.GetCompletedJobs(sessionkey)
			joblist.GetFailedJobs(sessionkey)
		} */

		for _, minion := range all_minions {
			if minion.JobID == 3 && minion.Migration_Stage == "waiting" {
				result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
				if result.Error != nil {
					logger.Errorf("failed to get minion %s from database", minion.Minion_Name)
					return
				}
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "completed")
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
				continue
			}

			if minion.JobID != 0 {
				status, err := Match_Job(sessionkey, minion, groupsdata)
				//stage_time := time.Now().String()
				//logger.Debugf(" %s: Minion %s JobID: %d, stage: %s at status: %s", stage_time, minion.Minion_Name, minion.JobID, minion.Migration_Stage, status)
				if err != nil {
					logger.Errorf("failed to get job status in Match_Job.")
					return
				}
				if status == "pending" {
					logger.Debugf("Minion %s Job %s is still in pending state.", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "pending")
					continue
				}
				if status == "completed" {
					logger.Debugf("Minion %s Job %s is completed.", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "completed")
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
					continue
				}

				if status == "failed" {
					logger.Infof("Minion %s Job %s is failed.", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "failed")
					db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "failed")
					continue
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
	logger.Infof("Check_Jobs final deadline reached. Exiting.")
	return
}

func Match_Job(sessionkey *auth.SumaSessionKey, minion Minion_Data, groupsdata *Update_Groups) (string, error) {

	if minion.JobID == 3 {
		//logger.Infof("Minion %s is not in any job. Maybe job is deleted. Set minion stage to completed.", minion.Minion_Name)
		return "completed", nil
	}

	status, err := Check_System_In_Jobs(sessionkey, minion.JobID, minion, groupsdata)
	if err != nil {
		logger.Errorln("failed to get job status in Check_System_In_Jobs.")
		return "", err
	}
	if status == "pending" {
		logger.Debugf("Minion %s is still in pending state.", minion.Minion_Name)
		return "pending", nil
	}
	if status == "completed" {
		logger.Debugf("Minion %s is completed.", minion.Minion_Name)
		return "completed", nil
	}
	if status == "failed" {
		logger.Debugf("Minion %s is failed.", minion.Minion_Name)
		return "failed", nil
	}

	if status == "not found" {
		if minion.JobID == 3 {
			//logger.Infof("Minion %s is not in any job. Maybe job is deleted. Set minion stage to completed.", minion.Minion_Name)
			return "completed", nil
		}
		logger.Infof("Minion %s is not in any job. Maybe job is deleted. Set minion stage to completed.", minion.Minion_Name)
		return "completed", nil
	}
	return "", nil
}
