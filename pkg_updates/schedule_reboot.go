package pkg_updates

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

type Schedule_Reboot_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Reboot_Response struct {
	JobID int `xmlrpc:"id"`
}

func Reboot(sessionkey *auth.SumaSessionKey, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	method := "system.scheduleReboot"

	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}

		//fmt.Printf("-----------Query DB Reboot %d\n", result.RowsAffected)
		//logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			if minion.JobID == 0 && minion.Migration_Stage == stage {
				logger.Debugf("Minion %s: set %s as completed due to manual intervention.\n", minion.Minion_Name, stage)
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "Completed")
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				continue
			}

			if minion.Target_Ident == "" && (minion.Migration_Stage == "spmigration_run" || minion.Migration_Stage == "spmigration_dryrun") {
				logger.Debugf("Target Ident is empty for minion %s\n", minion.Minion_Name)
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				/* subject := "Target Ident is empty"
				note := fmt.Sprintf("No valid migration target found. %s", minion.Minion_Name)
				Add_Note(sessionkey, minion.Minion_ID, subject, note) */
				continue
			}

			logger.Debugf("Minion %s starts %s stage.\n", minion.Minion_Name, stage)

			schedule_reboot_request := Schedule_Reboot_Request{
				Sessionkey:         sessionkey.Sessionkey,
				Sid:                minion.Minion_ID,
				EarliestOccurrence: time.Now(),
			}

			buf, err := gorillaxml.EncodeClientRequest(method, &schedule_reboot_request)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
			resp, err := request.MakeRequest(buf)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", string(buf))
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))

			/* responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatalf("ReadAll error: %s\n", err)
			}
			logger.Infof("responseBody: %s\n", responseBody) */
			reply := new(Schedule_Reboot_Response)
			err = gorillaxml.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				logger.Fatalf("Decode reboot Job response body failed: %s\n", err)
			}

			if reply.JobID > 0 {
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", reply.JobID)
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
				db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				logger.Infof("Minion %s has been scheduled to reboot\n", minion.Minion_Name)
			}
		}

	}
}
