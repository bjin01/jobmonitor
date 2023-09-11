package spmigration

import (
	"fmt"
	"log"

	"github.com/bjin01/jobmonitor/auth"
)

func (t *Target_Minions) Analyze_Pending_SPMigration(sessionkey *auth.SumaSessionKey,
	groupsdata *Migration_Groups, health *bool) {
	// get all minions which Migration Stage is in Product Migration and status is pending
	var analyze_target_minions Target_Minions
	analyze_target_minions.Tracking_file_name = fmt.Sprintf("%s.analyze", t.Tracking_file_name)
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
	analyze_target_minions.Write_Tracking_file()
	log.Printf("Execute analyze pending sp migration for %d minions\n", len(analyze_target_minions.Minion_List))
	analyze_target_minions.Reschedule_Pkg_Refresh(sessionkey)
	analyze_target_minions.Check_Pkg_Refresh_Jobs(sessionkey, health)
	jobid := analyze_target_minions.Schedule_Package_Updates(sessionkey)
	analyze_target_minions.Check_Package_Updates_Jobs(sessionkey, jobid, health)
	analyze_target_minions.Schedule_Reboot(sessionkey)
	analyze_target_minions.Check_Reboot_Jobs(sessionkey, health)
	analyze_target_minions.ListMigrationTarget(sessionkey, groupsdata)
	analyze_target_minions.Schedule_Migration(sessionkey, groupsdata, true)
	analyze_target_minions.Check_SP_Migration(sessionkey, true, health)
	analyze_target_minions.Schedule_Migration(sessionkey, groupsdata, false)
	analyze_target_minions.Check_SP_Migration(sessionkey, false, health)
	analyze_target_minions.Schedule_Reboot(sessionkey)
	analyze_target_minions.Check_Reboot_Jobs(sessionkey, health)
}

func (t *Target_Minions) Reschedule_Pkg_Refresh(sessionkey *auth.SumaSessionKey) {
	for i, recover_minion := range t.Minion_List {
		jobid, err := api_request_pkg_refresh(sessionkey, recover_minion.Minion_ID)
		if err != nil {
			log.Printf("api_request_pkg_refresh Error: %s\n", err)
		}
		if jobid > 0 {
			var host_info Host_Job_Info
			host_info.Pkg_Refresh_Job.JobID = jobid
			host_info.Pkg_Refresh_Job.JobStatus = "Scheduled"
			t.Minion_List[i].Host_Job_Info = host_info
			t.Minion_List[i].Migration_Stage = "Pkg_Refresh"
			t.Minion_List[i].Migration_Stage_Status = "Scheduled"
		} else {
			log.Printf("Minion %s is not ready for package refresh\n", recover_minion.Minion_Name)
			continue
		}
	}
}
