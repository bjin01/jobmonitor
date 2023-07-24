package spmigration

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Reboot_Jobs(sessionkey *auth.SumaSessionKey) {

	deadline := time.Now().Add(time.Duration(10) * time.Minute)

	for time.Now().Before(deadline) {
		var l schedules.ListJobs
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)
		l.GetPendingjobs(sessionkey)
		time.Sleep(10 * time.Second)
		t.Find_Reboot_Jobs(&l)

		if len(l.Pending.Result) == 0 {
			log.Printf("No more reboot job. Exit job check.\n")
			deadline = time.Now()
			break
		}

		log.Printf("Reboot Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			log.Printf("Reboot Job Status: %s %s\n", Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus,
				Minion.Minion_Name)

		}
		time.Sleep(10 * time.Second)
	}
	log.Printf("Reboot Job Status check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Reboot_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {

			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//fmt.Printf("Reboot Pending ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Pending"

			}
		}

		for _, p := range alljobs.Completed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//fmt.Printf("Reboot Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Completed"

			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//fmt.Printf("Reboot Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Failed"

			}
		}
	}
}