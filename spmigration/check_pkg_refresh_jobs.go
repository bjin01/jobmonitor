package spmigration

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Pkg_Refresh_Jobs(sessionkey *auth.SumaSessionKey, health *bool) {

	deadline := time.Now().Add(time.Duration(15) * time.Minute)
	extended_deadline_counter := 0

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
		t.Find_Pkg_Refresh_Jobs(&l)
		t.Find_Pkg_Refresh_Jobs_No_Targets(&l)

		if l.Found_Pending_Jobs == false {
			logger.Infof("No more pending pkg refresh job. Exit job check.\n")

			if extended_deadline_counter == 0 {
				deadline = time.Now().Add(time.Duration(240) * time.Second)
				extended_deadline_counter++
				continue
			}
			//break
		}

		logger.Infof("Package refresh Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			logger.Infof("Package refresh Job Status: %s %s %s\n", Minion.Migration_Stage,
				Minion.Migration_Stage_Status, Minion.Minion_Name)

		}

		for _, Minion := range t.No_Targets_Minions {
			logger.Infof("Package refresh Job Status: %s %s %s\n", Minion.Migration_Stage,
				Minion.Migration_Stage_Status, Minion.Minion_Name)

		}
		time.Sleep(10 * time.Second)
		t.Write_Tracking_file()
	}
	logger.Infof("Package refresh Job check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Pkg_Refresh_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("Pkg Refresh Pending Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//logger.Infof("Pkg Refresh Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Completed"

			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//logger.Infof("Pkg Refresh Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}

func (t *Target_Minions) Find_Pkg_Refresh_Jobs_No_Targets(alljobs *schedules.ListJobs) {
	for m, Minion := range t.No_Targets_Minions {
		for _, p := range alljobs.Pending.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("Pkg Refresh Pending Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Pending"
				t.No_Targets_Minions[m].Migration_Stage = "Pkg_Refresh"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//logger.Infof("Pkg Refresh Completed Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Completed"
				t.No_Targets_Minions[m].Migration_Stage = "Pkg_Refresh"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Completed"

			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//logger.Infof("Pkg Refresh Failed Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Failed"
				t.No_Targets_Minions[m].Migration_Stage = "Pkg_Refresh"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}
