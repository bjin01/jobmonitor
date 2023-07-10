package groups

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Assigne_Channels_Jobs(sessionkey *auth.SumaSessionKey) {

	deadline := time.Now().Add(time.Duration(2) * time.Minute)

	for time.Now().Before(deadline) {
		var l schedules.ListJobs
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)
		l.GetPendingjobs(sessionkey)
		time.Sleep(10 * time.Second)
		t.Find_Assigne_Channels_Jobs(&l)
		log.Printf("Assign Channels Job check 20 seconds. Deadline is %+v\n", deadline)
		for _, Minion := range t.Minion_List {
			log.Printf("Assign Channels Job Status: %s %s\n", Minion.Host_Job_Info.Assigne_Channels_Job.JobStatus,
				Minion.Minion_Name)

		}
		time.Sleep(10 * time.Second)
	}
	log.Printf("Assign Channels Job check deadline reached. %+v\n", deadline)
	return
}

func (t *Target_Minions) Find_Assigne_Channels_Jobs(alljobs *schedules.ListJobs) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {
			if p.Id == Minion.Host_Job_Info.Assigne_Channels_Job.JobID {
				//fmt.Printf("Pending Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Assigne_Channels_Job.JobStatus = "Pending"
				t.Minion_List[m].Migration_Stage = "Assign_Channels"
				t.Minion_List[m].Migration_Stage_Status = "Pending"

			}
		}

		for _, c := range alljobs.Completed.Result {
			if c.Id == Minion.Host_Job_Info.Assigne_Channels_Job.JobID {
				//fmt.Printf("Completed Job ID: %d\n", c.Id)
				t.Minion_List[m].Host_Job_Info.Assigne_Channels_Job.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Assign_Channels"
				t.Minion_List[m].Migration_Stage_Status = "Completed"

			}
		}

		for _, f := range alljobs.Failed.Result {
			if f.Id == Minion.Host_Job_Info.Assigne_Channels_Job.JobID {
				//fmt.Printf("Failed Job ID: %d\n", f.Id)
				t.Minion_List[m].Host_Job_Info.Assigne_Channels_Job.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Assign_Channels"
				t.Minion_List[m].Migration_Stage_Status = "Failed"

			}
		}
	}
}
