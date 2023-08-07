package spmigration

import (
	"fmt"
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ScheduleSPMigrationDryRun_Request struct {
	Sessionkey                                  string    `xmlrpc:"sessionKey"`
	Sid                                         int       `xmlrpc:"sid"`
	TargetIdent                                 string    `xmlrpc:"targetIdent"`
	BaseChannelLabel                            string    `xmlrpc:"baseChannelLabel"`
	OptionalChildChannels                       []string  `xmlrpc:"optionalChildChannels"`
	DryRun                                      bool      `xmlrpc:"dryRun"`
	AllowVendorChange                           bool      `xmlrpc:"allowVendorChange"`
	RemoveProductsWithNoSuccessorAfterMigration bool      `xmlrpc:"removeProductsWithNoSuccessorAfterMigration"`
	EarliestOccurrence                          time.Time `xmlrpc:"earliestOccurrence"`
}

type ScheduleSPMigrationDryRun_Response struct {
	JobID int `xmlrpc:"id"`
}

func (t *Target_Minions) Schedule_Migration(sessionkey *auth.SumaSessionKey,
	UserData *Migration_Groups, dryrun bool) {
	method := "system.scheduleProductMigration"

	for i, minion := range t.Minion_List {
		if dryrun == true && minion.Migration_Stage_Status == "Completed" &&
			minion.Migration_Stage == "Pkg_Refresh" {
			log.Printf("Minion %s is ready for migration dryrun\n", minion.Minion_Name)
		} else if dryrun == false && minion.Migration_Stage_Status == "Completed" &&
			minion.Migration_Stage == "Product Migration DryRun" {
			log.Printf("Minion %s is ready for migration\n", minion.Minion_Name)
		} else {
			continue
		}

		if minion.Target_Ident == "" {
			log.Default().Printf("Target Ident is empty for minion %s\n", minion.Minion_Name)
			continue
		}
		schedule_spmigration_request := ScheduleSPMigrationDryRun_Request{}
		schedule_spmigration_request.Sessionkey = sessionkey.Sessionkey
		schedule_spmigration_request.Sid = minion.Minion_ID
		schedule_spmigration_request.TargetIdent = minion.Target_Ident
		schedule_spmigration_request.BaseChannelLabel = minion.Target_base_channel
		schedule_spmigration_request.OptionalChildChannels = []string{}
		schedule_spmigration_request.AllowVendorChange = true
		schedule_spmigration_request.RemoveProductsWithNoSuccessorAfterMigration = true
		schedule_spmigration_request.EarliestOccurrence = time.Now()

		if dryrun == true {
			schedule_spmigration_request.DryRun = true
		} else {
			log.Printf("Schedule Product Migration for %s now!\n", minion.Minion_Name)
			schedule_spmigration_request.DryRun = false
		}

		buf, err := gorillaxml.EncodeClientRequest(method, &schedule_spmigration_request)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		fmt.Printf("client request spmigration buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Printf("Encoding scheduleProductMigration error: %s\n", err)
		}
		reply := new(ScheduleSPMigrationDryRun_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			if dryrun == true {
				log.Printf("Decode scheduleProductMigration_DryRun Job response body failed: %s %s\n", err, minion.Minion_Name)
			} else {
				log.Printf("Decode scheduleProductMigration Job response body failed: %s %s\n", err, minion.Minion_Name)
			}
		}
		//log.Printf("scheduleProductMigration_DryRun JobID: %d\n", reply.JobID)
		var host_info Host_Job_Info
		if dryrun == true {
			host_info.SP_Migration_DryRun_Job.JobID = reply.JobID
			host_info.SP_Migration_DryRun_Job.JobStatus = "Scheduled"
		} else {
			host_info.SP_Migration_Job.JobID = reply.JobID
			host_info.SP_Migration_Job.JobStatus = "Scheduled"
		}

		if reply.JobID > 0 {
			t.Minion_List[i].Host_Job_Info = host_info
			if dryrun == true {
				t.Minion_List[i].Migration_Stage = "Product Migration DryRun"
				log.Printf("Product migration dryrun JobID: %d %s\n",
					t.Minion_List[i].Host_Job_Info.SP_Migration_DryRun_Job.JobID, minion.Minion_Name)
			} else {
				t.Minion_List[i].Migration_Stage = "Product Migration"
				log.Printf("Product migration JobID: %d %s\n", t.Minion_List[i].Host_Job_Info.SP_Migration_Job.JobID, minion.Minion_Name)
			}
			t.Minion_List[i].Migration_Stage_Status = "Scheduled"

		} else {
			if dryrun == true {
				log.Printf("Minion %s product migration dryrun not possible.\n", minion.Minion_Name)
			} else {
				log.Printf("Minion %s product migration not possible.\n", minion.Minion_Name)
			}
			continue
		}
	}
	t.Write_Tracking_file()
}
