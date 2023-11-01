package spmigration

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/schedules"
)

func (t *Target_Minions) Check_Reboot_Jobs(sessionkey *auth.SumaSessionKey, email_job email.Job_Email_Body, jobinfo Email_job_info, health *bool) {

	deadline := time.Now().Add(time.Duration(t.Reboot_Timeout) * time.Minute)

	for time.Now().Before(deadline) {
		var l schedules.ListJobs
		if *health == false {
			logger.Infof("SPMigration can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		//logger.Infof("jobinfo in check_reboot_jobs.go: %+v\n", jobinfo)
		l.Found_Pending_Jobs = false
		l.GetPendingjobs(sessionkey)
		l.GetCompletedJobs(sessionkey)
		l.GetFailedJobs(sessionkey)

		time.Sleep(10 * time.Second)
		t.Find_Reboot_Jobs(&l, &email_job, &jobinfo)
		t.Find_Reboot_Jobs_No_Targets(&l, &email_job, &jobinfo)

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

func (t *Target_Minions) Find_Reboot_Jobs(alljobs *schedules.ListJobs, email_job *email.Job_Email_Body, jobinfo *Email_job_info) {
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
				/* //logger.Infof("Minion Name %s - Reboot Completed.\n", Minion.Minion_Name)
				email_job.Job_Response.Server_name = Minion.Minion_Name
				//logger.Infof("Minion ID %d - Reboot Completed.\n", Minion.Minion_ID)
				email_job.Job_Response.Server_id = Minion.Minion_ID
				//logger.Infof("Job ID %d - Reboot Completed.\n", p.Id)
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Pre Migration Reboot Completed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job) */
			} else if p.Id == Minion.Host_Job_Info.Reboot_Post_MigrationJob.JobID {
				//logger.Infof("Reboot Completed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Post_MigrationJob.JobStatus = "Completed"
				t.Minion_List[m].Migration_Stage = "Post Migration Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Completed"
				/* logger.Infof("Minion Name %s - Reboot Completed.\n", Minion.Minion_Name)
				email_job.Job_Response.Server_name = Minion.Minion_Name
				logger.Infof("Minion ID %d - Reboot Completed.\n", Minion.Minion_ID)
				email_job.Job_Response.Server_id = Minion.Minion_ID
				logger.Infof("Job ID %d - Reboot Completed.\n", p.Id)
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Pre Migration Reboot Completed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job) */
			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
				email_job.Job_Response.Server_name = Minion.Minion_Name
				email_job.Job_Response.Server_id = Minion.Minion_ID
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Pre Migration Reboot failed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job)
			} else if p.Id == Minion.Host_Job_Info.Reboot_Post_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.Minion_List[m].Host_Job_Info.Reboot_Post_MigrationJob.JobStatus = "Failed"
				t.Minion_List[m].Migration_Stage = "Post Migration Reboot"
				t.Minion_List[m].Migration_Stage_Status = "Failed"
				email_job.Job_Response.Server_name = Minion.Minion_Name
				email_job.Job_Response.Server_id = Minion.Minion_ID
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Post Migration Reboot failed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job)
			}
		}
	}
}

func (t *Target_Minions) Find_Reboot_Jobs_No_Targets(alljobs *schedules.ListJobs, email_job *email.Job_Email_Body, jobinfo *Email_job_info) {
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
				/* logger.Infof("Minion Name %s - Reboot Completed.\n", Minion.Minion_Name)
				email_job.Job_Response.Server_name = Minion.Minion_Name
				logger.Infof("Minion ID %d - Reboot Completed.\n", Minion.Minion_ID)
				email_job.Job_Response.Server_id = Minion.Minion_ID
				logger.Infof("Job ID %d - Reboot Completed.\n", p.Id)
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Reboot Completed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job) */
			}
		}

		for _, p := range alljobs.Failed.Result {
			if p.Id == Minion.Host_Job_Info.Reboot_Pre_MigrationJob.JobID {
				//logger.Infof("Reboot Failed Job ID: %d\n", p.Id)
				t.No_Targets_Minions[m].Host_Job_Info.Reboot_Pre_MigrationJob.JobStatus = "Failed"
				t.No_Targets_Minions[m].Migration_Stage = "Reboot"
				t.No_Targets_Minions[m].Migration_Stage_Status = "Failed"
				email_job.Job_Response.Server_name = Minion.Minion_Name
				email_job.Job_Response.Server_id = Minion.Minion_ID
				email_job.Job_Response.Job_ID = p.Id
				email_job.Job_Response.Job_Status = "Reboot failed"
				email_job.Job_Response.T7user = email_job.T7user
				jobinfo.Send_Job_Response_Email(*email_job)
			}
		}
	}
}
