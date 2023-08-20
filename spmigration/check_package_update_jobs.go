package spmigration

import (
	"fmt"
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ListSystemInJobs struct {
	ListInProgressSystems ListSystemInJobs_Response
	ListFailedSystems     ListSystemInJobs_Response
	ListCompletedSystems  ListSystemInJobs_Response
}

/* type ListSystemInJobs_Response struct {
	Result []struct {
		Server_name  string    `xmlrpc:"server_name,omitempty"`
		Base_channel string    `xmlrpc:"base_channel,omitempty"`
		Server_id    int       `xmlrpc:"server_id,omitempty"`
		Timestamp    time.Time `xmlrpc:"timestamp,omitempty"`
	}
} */

type ListSystemInJobs_Response struct {
	Result []struct {
		Server_name  string
		Base_channel string
		Server_id    int
		Timestamp    time.Time
	}
}

type ListSystemInJobs_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	ActionId   int    `xmlrpc:"actionId"`
}

func (t *Target_Minions) Check_Package_Updates_Jobs(sessionkey *auth.SumaSessionKey, jobid_pkg_update int, health *bool) {
	if jobid_pkg_update == 0 {
		log.Printf("No package update job scheduled. Exit check.\n")
		return
	}

	current_ListSystemInJobs_status := new(ListSystemInJobs)

	deadline := time.Now().Add(time.Duration(20) * time.Minute)

	for time.Now().Before(deadline) {
		if *health == false {
			log.Printf("SPMigration can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		log.Printf("Package Update Job loop check 20 seconds. Deadline is %+v\n", deadline)
		time.Sleep(10 * time.Second)
		current_ListSystemInJobs_status.List_InProgress_Systems(sessionkey, jobid_pkg_update)
		current_ListSystemInJobs_status.List_Failed_Systems(sessionkey, jobid_pkg_update)
		current_ListSystemInJobs_status.List_Completed_Systems(sessionkey, jobid_pkg_update)

		if len(current_ListSystemInJobs_status.ListInProgressSystems.Result) > 0 {
			log.Printf("Update Pkg bundle job ID: %d: ListInProgressSystems: %v\n", jobid_pkg_update,
				current_ListSystemInJobs_status.ListInProgressSystems)
			for i, minion := range t.Minion_List {
				for _, inprogress := range current_ListSystemInJobs_status.ListInProgressSystems.Result {
					if minion.Minion_ID == inprogress.Server_id {
						t.Minion_List[i].Migration_Stage = "Package Update"
						t.Minion_List[i].Migration_Stage_Status = "Pending"
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobID = jobid_pkg_update
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobStatus = "Pending"

					}
				}
			}
		} else {
			log.Printf("Update Pkg bundle job ID: %d: no more pending systems. Exit job check\n", jobid_pkg_update)
			deadline = time.Now()
		}

		if len(current_ListSystemInJobs_status.ListCompletedSystems.Result) > 0 {
			fmt.Printf("Update Pkg bundle job ID: %d: ListCompletedSystems: %v\n", jobid_pkg_update,
				current_ListSystemInJobs_status.ListCompletedSystems)
			for i, minion := range t.Minion_List {
				for _, completed := range current_ListSystemInJobs_status.ListCompletedSystems.Result {
					if minion.Minion_ID == completed.Server_id {
						t.Minion_List[i].Migration_Stage = "Package Update"
						t.Minion_List[i].Migration_Stage_Status = "Completed"
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobID = jobid_pkg_update
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobStatus = "Completed"

					}
				}
			}
		}
		if len(current_ListSystemInJobs_status.ListFailedSystems.Result) > 0 {
			fmt.Printf("Update Pkg bundle job ID: %d: ListFailedSystems: %v\n", jobid_pkg_update,
				current_ListSystemInJobs_status.ListFailedSystems)
			for i, minion := range t.Minion_List {
				for _, failed := range current_ListSystemInJobs_status.ListFailedSystems.Result {
					if minion.Minion_ID == failed.Server_id {
						t.Minion_List[i].Migration_Stage = "Package Update"
						t.Minion_List[i].Migration_Stage_Status = "Failed"
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobID = jobid_pkg_update
						t.Minion_List[i].Host_Job_Info.Update_Pkg_Job.JobStatus = "Failed"

					}
				}
			}
		}
		time.Sleep(10 * time.Second)
		t.Write_Tracking_file()
	}
	log.Printf("Package refresh Job check deadline reached. %+v\n", deadline)
	return
}

func (c *ListSystemInJobs) List_InProgress_Systems(sessionkey *auth.SumaSessionKey, jobid_pkg_update int) {
	request_obj := new(ListSystemInJobs_Request)

	request_obj.Sessionkey = sessionkey.Sessionkey
	request_obj.ActionId = jobid_pkg_update

	method := "schedule.listInProgressSystems"
	buf, err := gorillaxml.EncodeClientRequest(method, request_obj)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("request body: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		log.Fatalf("Decode ListSystemInJobs_Response Reponse body failed: %s\n", err)
	}

	c.ListInProgressSystems = *response_obj
}

func (c *ListSystemInJobs) List_Failed_Systems(sessionkey *auth.SumaSessionKey, jobid_pkg_update int) {
	request_obj := new(ListSystemInJobs_Request)

	request_obj.Sessionkey = sessionkey.Sessionkey
	request_obj.ActionId = jobid_pkg_update

	method := "schedule.listFailedSystems"
	buf, err := gorillaxml.EncodeClientRequest(method, request_obj)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("request body: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		log.Fatalf("Decode listFailedSystems Reponse body failed: %s\n", err)
	}

	c.ListFailedSystems = *response_obj
}

func (c *ListSystemInJobs) List_Completed_Systems(sessionkey *auth.SumaSessionKey, jobid_pkg_update int) {
	request_obj := new(ListSystemInJobs_Request)

	request_obj.Sessionkey = sessionkey.Sessionkey
	request_obj.ActionId = jobid_pkg_update

	method := "schedule.listCompletedSystems"
	buf, err := gorillaxml.EncodeClientRequest(method, request_obj)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("request body: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		log.Fatalf("Decode listCompletedSystems Reponse body failed: %s\n", err)
	}

	c.ListCompletedSystems = *response_obj
}
