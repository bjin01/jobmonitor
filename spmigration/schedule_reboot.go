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

	list_reboot_minions := []map[string]int{}

	for _, minion := range t.Minion_List {
		if minion.Migration_Stage_Status == "Completed" &&
			(minion.Migration_Stage == "Package Update" || minion.Migration_Stage == "Product Migration") {
			minion_values := map[string]int{minion.Minion_Name: minion.Minion_ID}
			list_reboot_minions = append(list_reboot_minions, minion_values)
			logger.Infof("Minion %s is ready for reboot\n", minion.Minion_Name)
		}
	}

	for _, minion := range t.No_Targets_Minions {
		if minion.Migration_Stage_Status == "Completed" &&
			(minion.Migration_Stage == "Package Update") {
			minion_values := map[string]int{minion.Minion_Name: minion.Minion_ID}
			list_reboot_minions = append(list_reboot_minions, minion_values)
			logger.Infof("Minion %s is ready for reboot\n", minion.Minion_Name)
		}
	}

	if len(list_reboot_minions) == 0 {
		logger.Infof("No minion is ready for reboot\n")
		return
	}

	for _, minion := range list_reboot_minions {
		sid := int(0)
		system_name := ""

		for key, value := range minion {
			sid = value
			system_name = key
		}
		schedule_reboot_request := Schedule_Reboot_Request{
			Sessionkey:         sessionkey.Sessionkey,
			Sid:                sid,
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
		logger.Infof("Reboot JobID: %d %s\n", reply.JobID, system_name)

		if reply.JobID > 0 {
			var host_info Host_Job_Info
			for i, exist_minion := range t.Minion_List {
				if exist_minion.Migration_Stage == "Package Update" {
					if exist_minion.Minion_ID == sid {
						host_info.Reboot_Pre_MigrationJob.JobID = reply.JobID
						host_info.Reboot_Pre_MigrationJob.JobStatus = "Scheduled"
						t.Minion_List[i].Host_Job_Info = host_info
						t.Minion_List[i].Migration_Stage = "Reboot"
						t.Minion_List[i].Migration_Stage_Status = "Scheduled"
					}
				}
				if exist_minion.Migration_Stage == "Product Migration" {
					if exist_minion.Minion_ID == sid {
						host_info.Reboot_Post_MigrationJob.JobID = reply.JobID
						host_info.Reboot_Post_MigrationJob.JobStatus = "Scheduled"
						t.Minion_List[i].Host_Job_Info = host_info
						t.Minion_List[i].Migration_Stage = "Post Migration Reboot"
						t.Minion_List[i].Migration_Stage_Status = "Scheduled"
					}
				}
			}

			for i, exist_minion := range t.No_Targets_Minions {
				if exist_minion.Migration_Stage == "Package Update" {
					if exist_minion.Minion_ID == sid {
						host_info.Reboot_Pre_MigrationJob.JobID = reply.JobID
						host_info.Reboot_Pre_MigrationJob.JobStatus = "Scheduled"
						t.No_Targets_Minions[i].Host_Job_Info = host_info
						t.No_Targets_Minions[i].Migration_Stage = "Reboot"
						t.No_Targets_Minions[i].Migration_Stage_Status = "Scheduled"
					}
				}
			}

		}

	}

	t.Write_Tracking_file()
}
