package spmigration

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_SP_Migration(sessionkey *auth.SumaSessionKey, dryrun bool, health *bool) {
	deadline := time.Now().Add(time.Duration(t.Jobcheck_Timeout) * time.Minute)
	if dryrun == true {
		logger.Infof("Dryrun mode. SP Migration DryRun Jobs will be monitored.\n")
		deadline = time.Now().Add(time.Duration(15) * time.Minute)
	}

	if dryrun == false {
		logger.Infof("SP Migration Jobs will be monitored.\n")
		deadline = time.Now().Add(time.Duration(t.Jobcheck_Timeout) * time.Minute)
	}

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
		t.Find_SPMigration_Jobs(&l, dryrun)

		if l.Found_Pending_Jobs == false {
			if dryrun == true {
				logger.Infof("No more pending spmigration dryrun job. Exit job check.\n")
			} else {
				logger.Infof("No more pending spmigration job. Exit job check.\n")
			}
			time.Sleep(60 * time.Second)
			deadline = time.Now()
			//break
		}

		if dryrun == true {
			logger.Infof("Spmigration dryrun Job check 20 seconds. Deadline is %+v\n", deadline)
		} else {
			logger.Infof("Spmigration Job check 20 seconds. Deadline is %+v\n", deadline)
		}

		for _, Minion := range t.Minion_List {
			if dryrun == true {
				logger.Infof("Spmigration dryrun Job Status: %s %s %s %d\n", Minion.Migration_Stage,
					Minion.Migration_Stage_Status, Minion.Minion_Name,
					Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID)
			} else {
				logger.Infof("Spmigration Job Status: %s %s %s %d\n", Minion.Migration_Stage,
					Minion.Migration_Stage_Status, Minion.Minion_Name, Minion.Host_Job_Info.SP_Migration_Job.JobID)
			}

		}
		time.Sleep(10 * time.Second)
		t.Write_Tracking_file()
	}
	if dryrun == true {
		logger.Infof("Spmigration dryrun Job check deadline reached. %+v\n", deadline)
	} else {
		logger.Infof("Spmigration Job check deadline reached. %+v\n", deadline)
	}
	return
}

func (t *Target_Minions) Find_SPMigration_Jobs(alljobs *schedules.ListJobs, dryrun bool) {
	for m, Minion := range t.Minion_List {
		for _, p := range alljobs.Pending.Result {
			jobid := new(int)
			if dryrun == true {
				if Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID != 0 {
					//log.Default().Printf("Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID: %d\n", Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID)
					*jobid = Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID
				} else {
					log.Default().Printf("NO: Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID found: %s.\n",
						Minion.Minion_Name)
				}
			} else {
				*jobid = Minion.Host_Job_Info.SP_Migration_Job.JobID
			}
			if p.Id == *jobid {
				alljobs.Found_Pending_Jobs = true
				//logger.Infof("SP Migration DryRun Pending Job ID: %d\n", p.Id)

				if dryrun == true {
					t.Minion_List[m].Host_Job_Info.SP_Migration_DryRun_Job.JobStatus = "Pending"
					t.Minion_List[m].Migration_Stage = "Product Migration DryRun"
					t.Minion_List[m].Migration_Stage_Status = "Pending"
				} else {
					t.Minion_List[m].Host_Job_Info.SP_Migration_Job.JobStatus = "Pending"
					t.Minion_List[m].Migration_Stage = "Product Migration"
					t.Minion_List[m].Migration_Stage_Status = "Pending"
				}
				//t.Minion_List[m].Migration_Stage_Status = "Pending"
			}
		}

		for _, p := range alljobs.Completed.Result {
			jobid := new(int)
			if dryrun == true {
				*jobid = Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID
			} else {
				*jobid = Minion.Host_Job_Info.SP_Migration_Job.JobID
			}

			if p.Id == *jobid {
				//logger.Infof("SP Migration DryRun Completed Job ID: %d\n", p.Id)

				if dryrun == true {
					t.Minion_List[m].Host_Job_Info.SP_Migration_DryRun_Job.JobStatus = "Completed"
					t.Minion_List[m].Migration_Stage = "Product Migration DryRun"
				} else {
					t.Minion_List[m].Host_Job_Info.SP_Migration_Job.JobStatus = "Completed"
					t.Minion_List[m].Migration_Stage = "Product Migration"
				}
				t.Minion_List[m].Migration_Stage_Status = "Completed"
			}
		}

		for _, p := range alljobs.Failed.Result {
			jobid := new(int)
			if dryrun == true {
				*jobid = Minion.Host_Job_Info.SP_Migration_DryRun_Job.JobID
			} else {
				*jobid = Minion.Host_Job_Info.SP_Migration_Job.JobID
			}
			if p.Id == *jobid {

				//logger.Infof("SP Migration DryRun Failed Job ID: %d\n", p.Id)

				if dryrun == true {
					t.Minion_List[m].Host_Job_Info.SP_Migration_DryRun_Job.JobStatus = "Failed"
					t.Minion_List[m].Migration_Stage = "Product Migration DryRun"
				} else {
					t.Minion_List[m].Host_Job_Info.SP_Migration_Job.JobStatus = "Failed"
					t.Minion_List[m].Migration_Stage = "Product Migration"
				}
				t.Minion_List[m].Migration_Stage_Status = "Failed"
			}
		}
	}
}
