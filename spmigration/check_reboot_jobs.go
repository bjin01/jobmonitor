package spmigration

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Reboot_Jobs(sessionkey *auth.SumaSessionKey, health *bool) {

	deadline := time.Now().Add(time.Duration(t.Reboot_Timeout) * time.Minute)

	for time.Now().Before(deadline) {
		var l schedules.ListJobs
		if *health == false {
			logger.Infof("SPMigration can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		l.Found_Pending_Jobs = false
		l.GetPendingjobs(sessionkey)
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)

		time.Sleep(10 * time.Second)
		t.Find_Reboot_Jobs(&l)
		t.Find_Reboot_Jobs_No_Targets(&l)

		if l.Found_Pending_Jobs == false {
			logger.Infof("No more reboot job. Exit job check.\n")
			deadline = time.Now()
			//break
		}

		logger.Infof("Reboot Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			if Minion.Migration_Stage == "Reboot" {
				logger.Infof("Reboot Job Status: %s %s %s\n", Minion.Migration_Stage, Minion.Migration_Stage_Status,
					Minion.Minion_Name)
			}

			if Minion.Migration_Stage == "Post Migration Reboot" {
				logger.Infof("Post Migration Reboot Job Status: %s %s %s\n", Minion.Migration_Stage, Minion.Migration_Stage_Status,
					Minion.Minion_Name)
			}
		}

		for _, Minion := range t.No_Targets_Minions {
			if Minion.Migration_Stage == "Reboot" {
				logger.Infof("Reboot Job Status: %s %s %s\n", Minion.Migration_Stage, Minion.Migration_Stage_Status,
					Minion.Minion_Name)
			}
		}
		time.Sleep(10 * time.Second)
		t.Write_Tracking_file()
	}
	logger.Infof("Reboot Job Status check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Reboot_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {

			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("Reboot Pending ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Pending"
			} else if p.Id == Minion.Host_Job_Info.Reboot_Post_MigrationJob.JobID {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("Reboot Pending ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Post_MigrationJob.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Post Migration Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Completed"
			} else if p.Id == Minion.Host_Job_Info.Reboot_Post_MigrationJob.JobID {
				//logger.Infof("Reboot Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Post_MigrationJob.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Post Migration Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Completed"
			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
			} else if p.Id == Minion.Host_Job_Info.Reboot_Post_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Post_MigrationJob.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Post Migration Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}

func (t *Target_Minions) Find_Reboot_Jobs_No_Targets(alljobs *schedules.ListJobs) {
	for m, Minion := range t.No_Targets_Minions {
		for _, p := range alljobs.Pending.Result {

			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("Reboot Pending ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Pending"
				t.No_Targets_Minions[m].Migration_Stage = "Reboot"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Completed Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Completed"
				t.No_Targets_Minions[m].Migration_Stage = "Reboot"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Completed"
			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Failed"
				t.No_Targets_Minions[m].Migration_Stage = "Reboot"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}
