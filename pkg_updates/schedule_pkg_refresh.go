package pkg_updates

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Pkg_Refresh_Response struct {
	JobID int `xmlrpc:"id"`
}

func Refresh_Packages(sessionkey *auth.SumaSessionKey, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	method := "system.schedulePackageRefresh"

	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}

		//fmt.Printf("-----------Query DB Pkg Refresh %d\n", result.RowsAffected)
		logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			if minion.JobID == 0 && minion.Migration_Stage == stage {
				logger.Infof("Minion %s: set reboot stage as completed due to manual intervention.\n", minion.Minion_Name)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "Completed")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				continue
			}
			logger.Infof("Minion %s starts %s stage.\n", minion.Minion_Name, stage)

			schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
				Sessionkey:         sessionkey.Sessionkey,
				Sid:                minion.Minion_ID,
				EarliestOccurrence: time.Now(),
			}

			buf, err := gorillaxml.EncodeClientRequest(method, &schedule_pkg_refresh_request)
			if err != nil {
				logger.Infof("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
			resp, err := request.MakeRequest(buf)
			if err != nil {
				logger.Infof("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", string(buf))
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))

			/* responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatalf("ReadAll error: %s\n", err)
			}
			logger.Infof("responseBody: %s\n", responseBody) */
			reply := new(Schedule_Pkg_Refresh_Response)
			err = gorillaxml.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				logger.Infof("Decode Pkg Refresh Job response body failed: %s\n", err)
			}
			logger.Infof("Package refresh JobID: %d\n", reply.JobID)
			var host_info Host_Job_Info
			host_info.Pkg_Refresh_Job.JobID = reply.JobID
			host_info.Pkg_Refresh_Job.JobStatus = "Scheduled"

			if reply.JobID > 0 {
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", reply.JobID)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				logger.Infof("Minion %s has been scheduled to package refresh.\n", minion.Minion_Name)
			}
		}
	} // end of for loop
}
