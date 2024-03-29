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
			logger.Infof("Minion %s is ready for migration dryrun\n", minion.Minion_Name)
		} else if dryrun == false && minion.Migration_Stage_Status == "Completed" &&
			minion.Migration_Stage == "Product Migration DryRun" {
			logger.Infof("Minion %s is ready for migration\n", minion.Minion_Name)
		} else {
			continue
		}

		if minion.Target_Ident == "" {
			log.Default().Printf("Target Ident is empty for minion %s\n", minion.Minion_Name)
			subject := "Target Ident is empty"
			note := fmt.Sprintf("No valid migration target found. %s", minion.Minion_Name)
			Add_Note(sessionkey, minion.Minion_ID, subject, note)
			continue
		}
		schedule_spmigration_request := ScheduleSPMigrationDryRun_Request{}
		schedule_spmigration_request.Sessionkey = sessionkey.Sessionkey
		schedule_spmigration_request.Sid = minion.Minion_ID
		schedule_spmigration_request.TargetIdent = minion.Target_Ident
		schedule_spmigration_request.BaseChannelLabel = minion.Target_base_channel
		schedule_spmigration_request.OptionalChildChannels = minion.Target_Optional_Channels
		schedule_spmigration_request.AllowVendorChange = true
		schedule_spmigration_request.RemoveProductsWithNoSuccessorAfterMigration = true
		schedule_spmigration_request.EarliestOccurrence = time.Now()

		/* for _, v := range UserData.Target_Products {
			if strings.TrimSpace(v.Product.Base_Channel) == minion.Target_base_channel {
				if len(v.Product.OptionalChildChannels) > 0 {
					for _, child := range v.Product.OptionalChildChannels {
						logger.Infof("%s: Add optional channel to schedule spmigration: %s\n",
							minion.Minion_Name, child)
						schedule_spmigration_request.OptionalChildChannels =
							append(schedule_spmigration_request.OptionalChildChannels, strings.TrimSpace(child))
					}
				}
			}
		} */

		if dryrun == true {
			schedule_spmigration_request.DryRun = true
		} else {
			logger.Infof("Schedule Product Migration for %s now!\n", minion.Minion_Name)
			schedule_spmigration_request.DryRun = false
		}

		buf, err := gorillaxml.EncodeClientRequest(method, &schedule_spmigration_request)
		if err != nil {
			logger.Fatalf("Encoding error: %s\n", err)
		}
		//logger.Infof("client request spmigration buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			logger.Infof("Encoding scheduleProductMigration error: %s\n", err)
		}
		//logger.Infof("scheduleProductMigration client request spmigration response: %s\n", resp.Body)
		defer resp.Body.Close()
		reply := new(ScheduleSPMigrationDryRun_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			if dryrun == true {
				logger.Infof("Decode scheduleProductMigration_DryRun Job response body failed: %s %s\n", err, minion.Minion_Name)
			} else {
				logger.Infof("Decode scheduleProductMigration Job response body failed: %s %s\n", err, minion.Minion_Name)
			}
		}
		//logger.Infof("scheduleProductMigration_DryRun JobID: %d\n", reply.JobID)
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
				logger.Infof("Product migration dryrun JobID: %d %s\n",
					t.Minion_List[i].Host_Job_Info.SP_Migration_DryRun_Job.JobID, minion.Minion_Name)
			} else {
				t.Minion_List[i].Migration_Stage = "Product Migration"
				logger.Infof("Product migration JobID: %d %s\n", t.Minion_List[i].Host_Job_Info.SP_Migration_Job.JobID, minion.Minion_Name)
			}
			t.Minion_List[i].Migration_Stage_Status = "Scheduled"

		} else {
			if dryrun == true {
				logger.Infof("Minion %s product migration dryrun not possible.\n", minion.Minion_Name)
			} else {
				logger.Infof("Minion %s product migration not possible.\n", minion.Minion_Name)
			}
			continue
		}
	}
	t.Write_Tracking_file()
}
