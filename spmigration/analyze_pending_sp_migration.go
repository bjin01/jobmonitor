package spmigration

import (
	"fmt"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
)

func (t *Target_Minions) Analyze_Pending_SPMigration(sessionkey *auth.SumaSessionKey,
	groupsdata *Migration_Groups, email_template_dir string, health *bool) {
	// get all minions which Migration Stage is in Product Migration and status is pending
	var analyze_target_minions Target_Minions
	emails := new(email.SPMigration_Email_Body)
	emails.Recipients = groupsdata.JobcheckerEmails
	email_job := new(email.Job_Email_Body)
	email_job.Recipients = groupsdata.JobcheckerEmails
	email_job.Template_dir = email_template_dir
	email_job.T7user = groupsdata.T7User

	jobinfo := new(Email_job_info)

	analyze_target_minions.Tracking_file_name = fmt.Sprintf("%s.analyse", t.Tracking_file_name)
	emails.SPmigration_Tracking_File = analyze_target_minions.Tracking_file_name
	emails.Template_dir = email_template_dir

	for _, minion := range t.Minion_List {
		if minion.Migration_Stage == "Product Migration" && minion.Migration_Stage_Status == "Pending" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Product Migration DryRun" && minion.Migration_Stage_Status == "Completed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Package Update" && minion.Migration_Stage_Status == "Pending" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Package Update" && minion.Migration_Stage_Status == "Completed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Pkg_Refresh" && minion.Migration_Stage_Status == "Failed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Pkg_Refresh" && minion.Migration_Stage_Status == "Pending" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Pkg_Refresh" && minion.Migration_Stage_Status == "Completed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Reboot" && minion.Migration_Stage_Status == "Pending" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Reboot" && minion.Migration_Stage_Status == "Failed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
		if minion.Migration_Stage == "Reboot" && minion.Migration_Stage_Status == "Completed" {
			analyze_target_minions.Minion_List = append(analyze_target_minions.Minion_List, minion)
		}
	}
	t.Write_Tracking_file()
	emails.Send_SPmigration_Email()
	logger.Infof("Execute analyze pending sp migration for %d minions\n", len(analyze_target_minions.Minion_List))

	analyze_target_minions.Reschedule_Pkg_Refresh(sessionkey)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Check_Pkg_Refresh_Jobs(sessionkey, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	jobid := analyze_target_minions.Schedule_Package_Updates(sessionkey)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Check_Package_Updates_Jobs(sessionkey, jobid, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Schedule_Reboot(sessionkey)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()
	analyze_target_minions.Check_Reboot_Jobs(sessionkey, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.ListMigrationTarget(sessionkey, groupsdata)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Schedule_Migration(sessionkey, groupsdata, true)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()
	analyze_target_minions.Check_SP_Migration(sessionkey, true, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Schedule_Migration(sessionkey, groupsdata, false)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()
	analyze_target_minions.Check_SP_Migration(sessionkey, false, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()

	analyze_target_minions.Salt_Set_Patch_Level(sessionkey, groupsdata)
	analyze_target_minions.Salt_Refresh_Grains(sessionkey, groupsdata)

	analyze_target_minions.Schedule_Reboot(sessionkey)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()
	analyze_target_minions.Check_Reboot_Jobs(sessionkey, *email_job, *jobinfo, health)
	t.Update_Target_Minion_Status(&analyze_target_minions)
	t.Write_Tracking_file()
}

func (t *Target_Minions) Reschedule_Pkg_Refresh(sessionkey *auth.SumaSessionKey) {
	for i, recover_minion := range t.Minion_List {
		jobid, err := api_request_pkg_refresh(sessionkey, recover_minion.Minion_ID)
		if err != nil {
			logger.Infof("api_request_pkg_refresh Error: %s\n", err)
		}
		if jobid > 0 {
			var host_info Host_Job_Info
			host_info.Pkg_Refresh_Job.JobID = jobid
			host_info.Pkg_Refresh_Job.JobStatus = "Scheduled"
			t.Minion_List[i].Host_Job_Info = host_info
			t.Minion_List[i].Migration_Stage = "Pkg_Refresh"
			t.Minion_List[i].Migration_Stage_Status = "Scheduled"
		} else {
			logger.Infof("Minion %s - scheduling package refresh failed.\n", recover_minion.Minion_Name)
			continue
		}
	}
}
