package spmigration

import (
	"fmt"
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
	minion_list := []map[string]int{}
	for _, minion := range t.Minion_List {
		minion_values := map[string]int{minion.Minion_Name: minion.Minion_ID}
		minion_list = append(minion_list, minion_values)
	}

	for _, minion := range t.No_Targets_Minions {
		minion_values := map[string]int{minion.Minion_Name: minion.Minion_ID}
		minion_list = append(minion_list, minion_values)
	}

	for _, minion := range minion_list {
		sid := int(0)
		system_name := ""

		for key, value := range minion {
			sid = value
			system_name = key
		}

		logger.Infof("Minion %s is ready for package refresh\n", system_name)

		schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
			Sessionkey:         sessionkey.Sessionkey,
			Sid:                sid,
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
			for i, exist_minion := range t.Minion_List {
				if exist_minion.Minion_ID == sid {
					t.Minion_List[i].Host_Job_Info = host_info
					t.Minion_List[i].Migration_Stage = "Pkgs_Refresh"
					t.Minion_List[i].Migration_Stage_Status = "Scheduled"
				}
			}
			for i, exist_minion := range t.No_Targets_Minions {
				if exist_minion.Minion_ID == sid {
					t.No_Targets_Minions[i].Host_Job_Info = host_info
					t.No_Targets_Minions[i].Migration_Stage = "Pkgs_Refresh"
					t.No_Targets_Minions[i].Migration_Stage_Status = "Scheduled"
				}
			}
		}
	} // end of for loop
}

func api_request_pkg_refresh(sessionkey *auth.SumaSessionKey, sid int) (int, error) {
	method := "system.schedulePackageRefresh"

	schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
		Sessionkey:         sessionkey.Sessionkey,
		Sid:                sid,
		EarliestOccurrence: time.Now(),
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &schedule_pkg_refresh_request)
	if err != nil {
		return 0, fmt.Errorf("Encoding error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		return 0, fmt.Errorf("Encoding error: %s\n", err)
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
		return 0, fmt.Errorf("Decode Pkg Refresh Job response body failed: %s\n", err)
	}
	logger.Infof("Package refresh JobID: %d\n", reply.JobID)
	return reply.JobID, nil
}
