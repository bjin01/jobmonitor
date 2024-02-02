package pkg_updates

import (
	"fmt"
	"strings"
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
		Message      string
	}
}

type ListSystemInJobs_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	ActionId   int    `xmlrpc:"actionId"`
}

func (s *ListSystemInJobs_Response) Get_Systems_In_Jobs() (Systems []string) {

	for _, result := range s.Result {
		Systems = append(Systems, result.Server_name)
	}
	return Systems
}

func Check_System_In_Jobs(sessionkey *auth.SumaSessionKey, jobid_pkg_update int, minion Minion_Data, groupsdata *Update_Groups) (string, error) {
	if jobid_pkg_update == 0 {
		logger.Infof("No Job ID provided. Exit check.")
		return "", fmt.Errorf("No Job ID provided. Exit check.")
	}

	current_ListSystemInJobs_status := new(ListSystemInJobs)
	current_ListSystemInJobs_status.List_InProgress_Systems(sessionkey, jobid_pkg_update)
	current_ListSystemInJobs_status.List_Failed_Systems(sessionkey, jobid_pkg_update)
	current_ListSystemInJobs_status.List_Completed_Systems(sessionkey, jobid_pkg_update)

	if len(current_ListSystemInJobs_status.ListInProgressSystems.Result) > 0 {
		logger.Debugf("Lookup job ID: %d: ListInProgressSystems: %v", jobid_pkg_update,
			current_ListSystemInJobs_status.ListInProgressSystems.Get_Systems_In_Jobs())

		for _, inprogress := range current_ListSystemInJobs_status.ListInProgressSystems.Result {
			if minion.Minion_ID == inprogress.Server_id {
				if strings.Contains(minion.Migration_Stage, "reboot") {
					go func() {
						//logger.Debugf("We do a reboot false-positive check here by using salt ping for %s.", minion.Minion_Name)
						Reboot_Triage(sessionkey, jobid_pkg_update, minion.Minion_ID, minion.Minion_Name, groupsdata)
					}()
				}
				return "pending", nil
			}
		}
	}

	if len(current_ListSystemInJobs_status.ListCompletedSystems.Result) > 0 {
		logger.Debugf("Lookup job ID: %d: ListCompletedSystems: %v", jobid_pkg_update,
			current_ListSystemInJobs_status.ListCompletedSystems.Get_Systems_In_Jobs())
		for _, completed := range current_ListSystemInJobs_status.ListCompletedSystems.Result {
			if minion.Minion_ID == completed.Server_id {
				return "completed", nil
			}
		}
	}

	if len(current_ListSystemInJobs_status.ListFailedSystems.Result) > 0 {
		logger.Debugf("Lookup job ID: %d: ListFailedSystems: %v", jobid_pkg_update,
			current_ListSystemInJobs_status.ListFailedSystems.Get_Systems_In_Jobs())

		for _, failed := range current_ListSystemInJobs_status.ListFailedSystems.Result {
			if minion.Minion_ID == failed.Server_id {
				return "failed", nil
			}
		}
	}
	return "not found", nil
}

func (c *ListSystemInJobs) List_InProgress_Systems(sessionkey *auth.SumaSessionKey, jobid_pkg_update int) {
	request_obj := new(ListSystemInJobs_Request)

	request_obj.Sessionkey = sessionkey.Sessionkey
	request_obj.ActionId = jobid_pkg_update

	method := "schedule.listInProgressSystems"
	buf, err := gorillaxml.EncodeClientRequest(method, request_obj)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}
	//logger.Infof("request body: %s", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s", err)
	}
	logger.Infof("responseBody: %s", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		logger.Fatalf("Decode ListSystemInJobs_Response Reponse body failed: %s", err)
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
		logger.Fatalf("Encoding error: %s", err)
	}
	//logger.Infof("request body: %s", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s", err)
	}
	logger.Infof("responseBody: %s", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		logger.Fatalf("Decode listFailedSystems Reponse body failed: %s", err)
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
		logger.Fatalf("Encoding error: %s", err)
	}
	//logger.Infof("request body: %s", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s", err)
	}
	logger.Infof("responseBody: %s", responseBody) */

	response_obj := new(ListSystemInJobs_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		logger.Fatalf("Decode listCompletedSystems Reponse body failed: %s", err)
	}

	c.ListCompletedSystems = *response_obj
}
