package spmigration

import (
	"log"
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
			log.Printf("SPMigration can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		l.Found_Pending_Jobs = false
		l.GetPendingjobs(sessionkey)
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)

		time.Sleep(10 * time.Second)
		t.Find_Pkg_Refresh_Jobs(&l)

		if l.Found_Pending_Jobs == false {
			log.Printf("No more pending pkg refresh job. Exit job check.\n")

			if extended_deadline_counter == 0 {
				deadline = time.Now().Add(time.Duration(40) * time.Second)
				extended_deadline_counter++
				continue
			} else {
				deadline = time.Now()
			}
			//break
		}

		log.Printf("Package refresh Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			log.Printf("Package refresh Job Status: %s %s %s\n", Minion.Migration_Stage,
				Minion.Migration_Stage_Status, Minion.Minion_Name)

		}
		time.Sleep(10 * time.Second)
		t.Write_Tracking_file()
	}
	log.Printf("Package refresh Job check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Pkg_Refresh_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				alljobs.Found_Pending_Jobs = true
				//fmt.Printf("Pkg Refresh Pending Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//fmt.Printf("Pkg Refresh Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Completed"

			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
				//fmt.Printf("Pkg Refresh Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Pkg_Refresh_Job.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Pkg_Refresh"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}
