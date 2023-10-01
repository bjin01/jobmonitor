package spmigration

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Schedule_Reboot_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Reboot_Response struct {
	JobID int `xmlrpc:"id"`
}

func (t *Target_Minions) Schedule_Reboot(sessionkey *auth.SumaSessionKey) {
	method := "system.scheduleReboot"

	for i, minion := range t.Minion_List {
		if minion.Migration_Stage_Status == "Completed" &&
			(minion.Migration_Stage == "Package Update" || minion.Migration_Stage == "Product Migration") {

			logger.Infof("Minion %s is ready for reboot\n", minion.Minion_Name)

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
			logger.Infof("Reboot JobID: %d %s\n", reply.JobID, minion.Minion_Name)

			if reply.JobID > 0 {
				var host_info Host_Job_Info
				if minion.Migration_Stage == "Package Update" {

					host_info.Reboot_Pre_MigrationJob.JobID = reply.JobID
					host_info.Reboot_Pre_MigrationJob.JobStatus = "Scheduled"
					t.Minion_List[i].Host_Job_Info = host_info
					t.Minion_List[i].Migration_Stage = "Reboot"
					t.Minion_List[i].Migration_Stage_Status = "Scheduled"
				} else if minion.Migration_Stage == "Product Migration" {

					host_info.Reboot_Post_MigrationJob.JobID = reply.JobID
					host_info.Reboot_Post_MigrationJob.JobStatus = "Scheduled"
					t.Minion_List[i].Host_Job_Info = host_info
					t.Minion_List[i].Migration_Stage = "Post Migration Reboot"
					t.Minion_List[i].Migration_Stage_Status = "Scheduled"
				} else {
					logger.Infof("Unknown Migration Stage: %s %s\n", minion.Migration_Stage, minion.Minion_Name)
				}
			} else {
				logger.Infof("Minion %s is not ready for reboot\n", minion.Minion_Name)
				continue
			}
		}
	}
	t.Write_Tracking_file()
}
