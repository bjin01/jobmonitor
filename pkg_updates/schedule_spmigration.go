package pkg_updates

import (
	"fmt"
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
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

func SPMigration(sessionkey *auth.SumaSessionKey, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string, dryrun bool) {

	method := "system.scheduleProductMigration"
	//get minion stage fromo DB

	//fmt.Printf("-----------Query DB Reboot %d\n", result.RowsAffected)

	for _, minion := range minion_list {
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}
		//logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			if minion.JobID == 0 && minion.Migration_Stage == stage {
				logger.Debugf("Minion %s: set %s as completed due to manual intervention.\n", minion.Minion_Name, stage)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "Completed")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				continue
			}

			logger.Debugf("Minion %s starts %s stage.\n", minion.Minion_Name, stage)

			if minion.Target_Ident == "" {
				log.Default().Printf("Target Ident is empty for minion %s\n", minion.Minion_Name)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				/* subject := "Target Ident is empty"
				note := fmt.Sprintf("No valid migration target found. %s", minion.Minion_Name)
				Add_Note(sessionkey, minion.Minion_ID, subject, note) */
				continue
			}

			var optional_channels []string
			for _, t := range minion.Target_Optional_Channels {
				if t.Channel_Label != "" {
					optional_channels = append(optional_channels, t.Channel_Label)
				}
			}

			fmt.Printf("----ident %s\n", minion.Target_Ident)
			fmt.Printf("-----base channel %s\n", minion.Target_base_channel)
			fmt.Printf("------Minion %s: optional channels: %s\n", minion.Minion_Name, optional_channels)
			schedule_spmigration_request := ScheduleSPMigrationDryRun_Request{}
			schedule_spmigration_request.Sessionkey = sessionkey.Sessionkey
			schedule_spmigration_request.Sid = minion.Minion_ID
			schedule_spmigration_request.TargetIdent = minion.Target_Ident
			schedule_spmigration_request.BaseChannelLabel = minion.Target_base_channel
			schedule_spmigration_request.OptionalChildChannels = optional_channels
			schedule_spmigration_request.AllowVendorChange = true
			schedule_spmigration_request.RemoveProductsWithNoSuccessorAfterMigration = true
			schedule_spmigration_request.EarliestOccurrence = time.Now()

			if dryrun == true {
				schedule_spmigration_request.DryRun = true
			} else {
				logger.Debugf("Schedule Product Migration for %s now!\n", minion.Minion_Name)
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

			if dryrun == true && reply.JobID > 0 {
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", reply.JobID)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				logger.Infof("Minion %s has been scheduled for spmigration dryrun.\n", minion.Minion_Name)
			}

			if dryrun == false && reply.JobID > 0 {
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", reply.JobID)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				logger.Infof("Minion %s has been scheduled for spmigration.\n", minion.Minion_Name)
			}
		}
	}

}
