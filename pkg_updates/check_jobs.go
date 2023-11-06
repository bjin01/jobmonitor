package pkg_updates

import (
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

func Check_Jobs(sessionkey *auth.SumaSessionKey, health *bool, db *gorm.DB) {

	deadline := time.Now().Add(time.Duration(60) * time.Minute)

	for time.Now().Before(deadline) {
		if *health == false {
			logger.Infof("Check_Jobs can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
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
			if minion.JobID == 3 {
				result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
				if result.Error != nil {
					logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
					return
				}
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "completed")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
				continue
			}

			if minion.JobID != 0 {
				status, err := Match_Job(sessionkey, minion)
				if err != nil {
					logger.Errorf("failed to get job status in Match_Job.")
					return
				}
				if status == "pending" {
					logger.Infof("Minion %s Job %s is still in pending state.\n", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "pending")
					continue
				}
				if status == "completed" {
					logger.Infof("Minion %s Job %s is completed.\n", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "completed")
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
					continue
				}

				if status == "failed" {
					logger.Infof("Minion %s Job %s is failed.\n", minion.Minion_Name, minion.Migration_Stage)
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "failed")
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "failed")
					continue
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
	logger.Infof("Check_Jobs final deadline reached. Exiting.\n")
	return
}

func Match_Job(sessionkey *auth.SumaSessionKey, minion Minion_Data) (string, error) {

	status, err := Check_System_In_Jobs(sessionkey, minion.JobID, minion)
	if err != nil {
		logger.Errorf("failed to get job status in Check_System_In_Jobs.")
		return "", err
	}
	if status == "pending" {
		logger.Infof("Minion %s is still in pending state.\n", minion.Minion_Name)
		return "pending", nil
	}
	if status == "completed" {
		logger.Infof("Minion %s is completed.\n", minion.Minion_Name)
		return "completed", nil
	}
	if status == "failed" {
		logger.Infof("Minion %s is failed.\n", minion.Minion_Name)
		return "failed", nil
	}

	if status == "not found" {
		logger.Infof("Minion %s is not in any job. Maybe job is deleted. Set minion stage to completed.\n", minion.Minion_Name)
		return "completed", nil
	}
	return "", nil
}
