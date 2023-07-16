package spmigration

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Pkg_Refresh_Jobs(sessionkey *auth.SumaSessionKey) {

	deadline := time.Now().Add(time.Duration(2) * time.Minute)

	for time.Now().Before(deadline) {
		var l schedules.ListJobs
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)
		l.GetPendingjobs(sessionkey)
		time.Sleep(10 * time.Second)
		t.Find_Pkg_Refresh_Jobs(&l)

		if len(l.Pending.Result) == 0 {
			log.Printf("No more pending pkg refresh job. Exit job check.\n")
			deadline = time.Now()
			break
		}

		log.Printf("Package refresh Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			log.Printf("Package refresh Job Status: %s %s\n", Minion.Host_Job_Info.Pkg_Refresh_Job.JobStatus,
				Minion.Minion_Name)

		}
		time.Sleep(10 * time.Second)
	}
	log.Printf("Package refresh Job check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Pkg_Refresh_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {
			if p.Id == Minion.Host_Job_Info.Pkg_Refresh_Job.JobID {
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
