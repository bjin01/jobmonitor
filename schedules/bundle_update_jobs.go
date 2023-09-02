package schedules

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

func (t *Jobstatus) Check_Package_Updates_Jobs(sessionkey *auth.SumaSessionKey, scheduled_jobs_by_minions []Job, jobid_pkg_update int, deadline time.Time) {
	current_ListSystemInJobs_status := new(ListSystemInJobs)

	current_ListSystemInJobs_status.List_InProgress_Systems(sessionkey, jobid_pkg_update)
	current_ListSystemInJobs_status.List_Failed_Systems(sessionkey, jobid_pkg_update)
	current_ListSystemInJobs_status.List_Completed_Systems(sessionkey, jobid_pkg_update)

	if len(current_ListSystemInJobs_status.ListInProgressSystems.Result) > 0 {
		//log.Printf("Update Pkg bundle job ID: %d: ListInProgressSystems: %v\n", jobid_pkg_update,
		//	current_ListSystemInJobs_status.ListInProgressSystems)
		for _, minion := range scheduled_jobs_by_minions {
			for _, inprogress := range current_ListSystemInJobs_status.ListInProgressSystems.Result {
				if minion.Hostname == inprogress.Server_name {
					minion.ServerID = inprogress.Server_id
					t.Pending = append(t.Pending, minion)
					log.Printf("Update Pkg bundle job ID: %d: Pending: %v\n", jobid_pkg_update, inprogress.Server_name)
				}
			}
		}
	} else {
		log.Printf("Update Pkg bundle job ID: %d: no more pending systems. Exit job check\n", jobid_pkg_update)
	}

	if len(current_ListSystemInJobs_status.ListCompletedSystems.Result) > 0 {
		//log.Printf("Update Pkg bundle job ID: %d: ListCompletedSystems: %v\n", jobid_pkg_update,
		//	current_ListSystemInJobs_status.ListCompletedSystems)
		for _, minion := range scheduled_jobs_by_minions {
			for _, completed := range current_ListSystemInJobs_status.ListCompletedSystems.Result {
				if minion.Hostname == completed.Server_name {
					minion.ServerID = completed.Server_id
					if len(current_ListSystemInJobs_status.ListInProgressSystems.Result) == 0 {
						err := Create_pkg_refresh_job(sessionkey, completed.Server_id, completed.Server_name)
						if err != nil {
							log.Printf("create_pkg_refresh_job error: %s\n", err)
						}
						log.Println("Sleep 120 seconds to allow package refresh job to complete")
						time.Sleep(120 * time.Second)
					}

					t.Completed = append(t.Completed, minion)
					log.Printf("Update Pkg bundle job ID: %d: Completed: %v\n", jobid_pkg_update, completed.Server_name)
				}
			}
		}
	}
	if len(current_ListSystemInJobs_status.ListFailedSystems.Result) > 0 {
		//log.Printf("Update Pkg bundle job ID: %d: ListFailedSystems: %v\n", jobid_pkg_update,
		//	current_ListSystemInJobs_status.ListFailedSystems)
		for _, minion := range scheduled_jobs_by_minions {
			for _, failed := range current_ListSystemInJobs_status.ListFailedSystems.Result {
				if minion.Hostname == failed.Server_name {
					minion.ServerID = failed.Server_id
					t.Failed = append(t.Failed, minion)
					log.Printf("Update Pkg bundle job ID: %d: Failed: %v\n", jobid_pkg_update, failed.Server_name)
				}
			}
		}
	}

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
