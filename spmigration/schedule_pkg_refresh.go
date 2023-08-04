package spmigration

import (
	"fmt"
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Pkg_Refresh_Response struct {
	JobID int `xmlrpc:"id"`
}

func (t *Target_Minions) Schedule_Pkg_refresh(sessionkey *auth.SumaSessionKey) {
	method := "system.schedulePackageRefresh"

	for i, minion := range t.Minion_List {
		if minion.Migration_Stage_Status == "Completed" && minion.Migration_Stage == "Reboot" {

			fmt.Printf("Minion %s is ready for package refresh\n", minion.Minion_Name)

			schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
				Sessionkey:         sessionkey.Sessionkey,
				Sid:                minion.Minion_ID,
				EarliestOccurrence: time.Now(),
			}

			buf, err := gorillaxml.EncodeClientRequest(method, &schedule_pkg_refresh_request)
			if err != nil {
				log.Fatalf("Encoding error: %s\n", err)
			}
			//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
			resp, err := request.MakeRequest(buf)
			if err != nil {
				log.Fatalf("Encoding error: %s\n", err)
			}
			//fmt.Printf("buffer: %s\n", string(buf))
			//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))

			/* responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("ReadAll error: %s\n", err)
			}
			fmt.Printf("responseBody: %s\n", responseBody) */
			reply := new(Schedule_Pkg_Refresh_Response)
			err = gorillaxml.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				log.Fatalf("Decode Pkg Refresh Job response body failed: %s\n", err)
			}
			log.Printf("Package refresh JobID: %d\n", reply.JobID)
			var host_info Host_Job_Info
			host_info.Pkg_Refresh_Job.JobID = reply.JobID
			host_info.Pkg_Refresh_Job.JobStatus = "Scheduled"

			if reply.JobID > 0 {
				t.Minion_List[i].Host_Job_Info = host_info
				t.Minion_List[i].Migration_Stage = "Pkgs_Refresh"
				t.Minion_List[i].Migration_Stage_Status = "Scheduled"
			}
		} else {
			log.Printf("Minion %s is not ready for package refresh\n", minion.Minion_Name)
			continue
		}
	}
}
