package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/request"
	"github.com/bjin01/jobmonitor/schedules"
)

func Jobmonitor(SUMAConfig *SUMAConfig, alljobs schedules.ScheduledJobs,
	instance_jobs_patching schedules.Jobs_Patching, templates_dir *email.Templates_Dir, health *bool) {
	//key := "R2bfp223Qsk-pX970Jw8tyJUChT4-e2J8anZ4G4n4IM="
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
	}

	var sumaconf Sumaconf
	for a, b := range SUMAConfig.SUMA {
		sumaconf.Server = a
		b.Password = Decrypt(key, b.Password)
		sumaconf.Password = b.Password
		sumaconf.User = b.User
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		log.Fatal(err)
	}

	jobstatus_result := new(schedules.Jobstatus)

	switch alljobs.JobType {
	case "patching":
		jobstatus_result.JobType = "patching"
	case "reboot":
		jobstatus_result.JobType = "reboot"
	default:
		jobstatus_result.JobType = "patching"
	}

	if len(instance_jobs_patching.JobcheckerEmails) != 0 {
		jobstatus_result.JobcheckerEmails = instance_jobs_patching.JobcheckerEmails
		log.Printf("JobcheckerEmails : %v\n", jobstatus_result.JobcheckerEmails)
	} else {
		jobstatus_result.JobcheckerEmails = []string{}
		log.Printf("JobcheckerEmails : %v\n", jobstatus_result.JobcheckerEmails)
	}

	if instance_jobs_patching.JobcheckerTimeout != 0 {
		jobstatus_result.JobcheckerTimeout = instance_jobs_patching.JobcheckerTimeout
		log.Printf("JobcheckerTimeout : %v\n", jobstatus_result.JobcheckerTimeout)
	} else {
		jobstatus_result.JobcheckerTimeout = 2
		log.Printf("JobcheckerTimeout : %v\n", jobstatus_result.JobcheckerTimeout)
	}

	if len(instance_jobs_patching.Offline_minions) != 0 {
		jobstatus_result.Offline_minions = instance_jobs_patching.Offline_minions
		log.Printf("Offline_minions : %v\n", jobstatus_result.Offline_minions)
	} else {
		jobstatus_result.Offline_minions = []string{}
		log.Printf("Offline_minions : %v\n", jobstatus_result.Offline_minions)
	}

	if len(instance_jobs_patching.Disqualified_minions) != 0 {
		jobstatus_result.Disqualified_minions = instance_jobs_patching.Disqualified_minions
		log.Printf("Disqualified_minions : %v\n", jobstatus_result.Disqualified_minions)
	} else {
		jobstatus_result.Disqualified_minions = []string{}
		log.Printf("Disqualified_minions : %v\n", jobstatus_result.Disqualified_minions)
	}

	if len(instance_jobs_patching.No_patch_execptions) != 0 {
		jobstatus_result.No_patch_execptions = instance_jobs_patching.No_patch_execptions
		log.Printf("No_patch_execptions : %v\n", jobstatus_result.No_patch_execptions)
	} else {
		jobstatus_result.No_patch_execptions = []string{}
		log.Printf("No_patch_execptions : %v\n", jobstatus_result.No_patch_execptions)
	}

	if instance_jobs_patching.JobstartDelay != 0 {
		jobstatus_result.JobstartDelay = instance_jobs_patching.JobstartDelay
		log.Printf("JobstartDelay : %v\n", jobstatus_result.JobstartDelay)
	} else {
		jobstatus_result.JobstartDelay = 1
		log.Printf("JobstartDelay : %v\n", jobstatus_result.JobstartDelay)
	}

	if instance_jobs_patching.T7user != "" {
		jobstatus_result.T7user = instance_jobs_patching.T7user
		log.Printf("T7user : %v\n", jobstatus_result.T7user)
	} else {
		jobstatus_result.T7user = "unknown"
		log.Printf("T7user : %v\n", jobstatus_result.T7user)
	}

	if instance_jobs_patching.Post_patching_file != "" {
		jobstatus_result.Post_patching_file = instance_jobs_patching.Post_patching_file
		log.Printf("Post_patching_file : %v\n", jobstatus_result.Post_patching_file)
	} else {
		jobstatus_result.Post_patching_file = ""
		log.Printf("Post_patching_file : %v\n", jobstatus_result.Post_patching_file)
	}

	if instance_jobs_patching.Post_patching != "" {
		jobstatus_result.Post_patching = instance_jobs_patching.Post_patching
		log.Printf("Post_patching state : %v\n", jobstatus_result.Post_patching)
	} else {
		jobstatus_result.Post_patching = ""
		log.Printf("Post_patching state : %v\n", jobstatus_result.Post_patching)
	}

	if instance_jobs_patching.Prep_patching != "" {
		jobstatus_result.Prep_patching = instance_jobs_patching.Prep_patching
		log.Printf("Prep_patching state : %v\n", jobstatus_result.Prep_patching)
	} else {
		jobstatus_result.Prep_patching = ""
		log.Printf("Prep_patching state : %v\n", jobstatus_result.Prep_patching)
	}

	if instance_jobs_patching.Patch_level != "" {
		jobstatus_result.Patch_level = instance_jobs_patching.Patch_level
		log.Printf("Patch_level : %v\n", jobstatus_result.Patch_level)
	} else {
		jobstatus_result.Patch_level = ""
		log.Printf("Patch_level : %v\n", jobstatus_result.Patch_level)
	}

	for _, job := range alljobs.AllJobs {
		log.Printf("Host: %s \tJob-ID: %d\n", job.Hostname, job.JobID)
	}

	if len(alljobs.AllJobs) != 0 {

		deadline := time.Now().Add(time.Duration(jobstatus_result.JobcheckerTimeout) * time.Minute)
		Jobstart_starttime := time.Now().Add(time.Duration(jobstatus_result.JobstartDelay) * time.Minute)
		jobstatus_result.JobStartTime = Jobstart_starttime.Format(time.RFC822Z)
		jobstatus_result.YamlFileName = fmt.Sprintf("completed_%s_%s", jobstatus_result.T7user, Jobstart_starttime.Format("20060102150405"))
	begin:
		for time.Now().Before(deadline) {
			if *health == false {
				log.Printf("SUMA Health check failed. Skip jobcheck loop. We will continue if SUMA is online again.")
				time.Sleep(20 * time.Second)
				continue
			}

			log.Printf("Looping every minute. Deadline is %+v\n", deadline)
			log.Printf("Jobcheck will start at %+v\n", Jobstart_starttime)
			for time.Now().After(Jobstart_starttime) {
				jobstatus_result.Pending = []schedules.Job{}
				jobstatus_result.Failed = []schedules.Job{}
				jobstatus_result.Completed = []schedules.Job{}
				jobstatus_result.Cancelled = []schedules.Job{}

				if alljobs.Full_Update_Jobs.Full_Update_Job_ID != nil {
					log.Printf("Monitor Full Update Job ID: %v\n", alljobs.Full_Update_Jobs.Full_Update_Job_ID)
					if len(alljobs.Full_Update_Jobs.Full_Update_Job_ID) > 0 {
						for _, j := range alljobs.Full_Update_Jobs.Full_Update_Job_ID {
							jobstatus_result.Check_Package_Updates_Jobs(SessionKey, alljobs.AllJobs, j)
						}
					}
				} else {
					jobstatus_result.Compare(SessionKey, alljobs.AllJobs)
				}

				if len(jobstatus_result.Pending) > 0 {
					log.Printf("Pending Jobs: %+v\n", jobstatus_result.Pending)
				}
				if len(jobstatus_result.Failed) > 0 {
					log.Printf("Failed Jobs: %+v\n", jobstatus_result.Failed)
				}
				if len(jobstatus_result.Completed) > 0 {
					log.Printf("Completed Jobs: %+v\n", jobstatus_result.Completed)
				}
				//log.Printf("Jobstatus result Pending: %+v\n", jobstatus_result.Pending)
				time.Sleep(30 * time.Second)
				if len(jobstatus_result.Pending) == 0 {
					log.Printf("No more pending Jobs. Exit loop. Email sent.")
					if jobstatus_result.JobType == "patching" {
						jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
						if err != nil {
							log.Default().Printf("ERROR: reboot list: %s\n", err)
						}
						jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result, templates_dir)
						if err != nil {
							log.Default().Printf("ERROR: reboot sls: %s\n", err)
						}
					}

					email.Sendit(jobstatus_result, templates_dir)
					break begin
				}
			}
			time.Sleep(60 * time.Second)
		}
		if len(jobstatus_result.Pending) > 0 {
			if jobstatus_result.JobType == "patching" {
				jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
				if err != nil {
					log.Default().Printf("ERROR: reboot list: %s\n", err)
				}
				jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result, templates_dir)
				if err != nil {
					log.Default().Printf("ERROR: reboot sls: %s\n", err)
				}
			}
			email.Sendit(jobstatus_result, templates_dir)
		}
	} else {
		log.Println("No Patch Jobs found.")
	}

	log.Println("Jobchecker timeout reached or not more jobs in pending status. Email sent.")
	err = auth.Logout("auth.logout", *SessionKey)
	if err != nil {
		log.Fatal(err)
	}
}
