package pkg_updates

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	"github.com/bjin01/jobmonitor/saltapi"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ListSystemEvents_Request struct {
	Sessionkey   string    `xmlrpc:"sessionKey"`
	Sid          int       `xmlrpc:"sid"`
	EarliestDate time.Time `xmlrpc:"earliest_date"`
}

type ListSystemEvents_Response struct {
	Result []struct {
		Failed_count     int       `xmlrpc:"failed_count,omitempty"`
		Modified_date    time.Time `xmlrpc:"modified_date,omitempty"`
		Action_type      string    `xmlrpc:"action_type,omitempty"`
		Created_date     time.Time `xmlrpc:"created_date,omitempty"`
		Successful_count int       `xmlrpc:"successful_count,omitempty"`
		Earliest_action  string    `xmlrpc:"earliest_action,omitempty"`
		Archived         int       `xmlrpc:"archived,omitempty"`
		Scheduler_user   string    `xmlrpc:"scheduler_user,omitempty"`
		Name             string    `xmlrpc:"name,omitempty"`
		Id               int       `xmlrpc:"id,omitempty"`
		Version          string    `xmlrpc:"version,omitempty"`
		Completed_date   time.Time `xmlrpc:"completed_date,omitempty"`
		Pickup_date      time.Time `xmlrpc:"pickup_date,omitempty"`
		Result_msg       string    `xmlrpc:"result_msg,omitempty"`
	}
}

func Reboot_Triage(sessionkey *auth.SumaSessionKey, jobid int, minion_id int, minion_name string, groupsdata *Update_Groups) {

	if CPU_Load() > 10 {
		logger.Infof("CPU load is too high, %d, Reboot_Triage will try later again.", CPU_Load())
		return
	}

	request_obj := new(ListSystemEvents_Request)

	request_obj.Sessionkey = sessionkey.Sessionkey
	request_obj.Sid = minion_id
	request_obj.EarliestDate = time.Now().Add(-30 * time.Minute)

	method := "system.listSystemEvents"

	buf, err := gorillaxml.EncodeClientRequest(method, request_obj)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}
	//logger.Infof("request body: %s from Reboot_Triage", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s", err)
	}

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s", err)
	}
	logger.Infof("responseBody of system.listSystemEvents: %s", responseBody) */

	response_obj := new(ListSystemEvents_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, response_obj)
	if err != nil {
		logger.Fatalf("Decode system.listSystemEvents Reponse body failed: %s", err)
	}
	/* logger.Infof("response_obj of system.listSystemEvents: %d", len(response_obj.Result))
	logger.Infof("response_obj of system.listSystemEvents: %v", response_obj.Result) */
	if len(response_obj.Result) > 0 {
		for _, event := range response_obj.Result {
			if event.Id == jobid {
				Reboot_triage_timer := 5
				if groupsdata.Reboot_Triage_Timer == 0 {
					logger.Debugf("The reboot_triage_timer is not set, set to default 15 minutes.")
					Reboot_triage_timer = 15
				}
				if groupsdata.Reboot_Triage_Timer < 15 {
					logger.Debugf("The reboot_triage_timer is too short - %d minutes, set to default 15 minutes.", groupsdata.Reboot_Triage_Timer)
					Reboot_triage_timer = 15
				} else {
					Reboot_triage_timer = groupsdata.Reboot_Triage_Timer
				}
				if time.Now().Before(event.Pickup_date.Add(time.Duration(Reboot_triage_timer) * time.Minute)) {
					logger.Debugf("Reboot_Triage() %s event ID %d by User: %s, picked up: %s, less than %d minutes, so not continue here.", minion_name, jobid, event.Scheduler_user, fmt.Sprintf("%s", event.Pickup_date), Reboot_triage_timer)
					continue
				}
				if !event.Pickup_date.IsZero() && event.Completed_date.IsZero() {
					logger.Infof("Reboot_Triage() %s event ID %d by User: %s, picked up: %s", minion_name, jobid, event.Scheduler_user, fmt.Sprintf("%s", event.Pickup_date))
					saltdata := new(saltapi.Salt_Data)
					saltdata.SaltMaster = groupsdata.SaltMaster_Address
					saltdata.SaltApi_Port = groupsdata.SaltApi_Port
					saltdata.Username = groupsdata.SaltUser
					saltdata.Password = groupsdata.SaltPassword
					saltdata.Target_List = []string{minion_name}
					saltdata.Arg = []string{""}
					if groupsdata.Timeout > 0 && groupsdata.GatherJobTimeout > 0 {
						saltdata.Arg = append(saltdata.Arg, fmt.Sprintf("timeout=%d", groupsdata.Timeout))
						saltdata.Arg = append(saltdata.Arg, fmt.Sprintf("gather_job_timeout=%d", groupsdata.GatherJobTimeout))
					}

					saltdata.Target_List = []string{minion_name}
					saltdata.Login()
					if saltdata.Run_Manage_Status_with_Response() {
						logger.Infof("Reboot_Triage() %s event ID %d by User: %s, minion is online", minion_name, jobid, event.Scheduler_user)
						saltdata.SaltCmd = "event.send"
						saltdata.Arg = []string{fmt.Sprintf("salt/minion/%s/start", minion_name)}
						logger.Infof("Reboot_Triage() %s event ID %d by User: %s, sending event: %s", minion_name, jobid, event.Scheduler_user, fmt.Sprintf("%s", saltdata.Arg))
						saltdata.Execute_Command_Async()
						continue
					} else {
						logger.Infof("Reboot_Triage() %s event ID %d by User: %s, picked up: %s, minion is still offline", minion_name, jobid, event.Scheduler_user, fmt.Sprintf("%s", event.Pickup_date))

					}

				} else {
					logger.Infof("Reboot_Triage() %s event ID %d by User: %s:%s, completed: %s", minion_name, jobid, event.Scheduler_user, event.Result_msg, fmt.Sprintf("%s", event.Completed_date))
				}
			}
		}
	}

}

func CPU_Load() int {
	out, err := exec.Command("cat", "/proc/loadavg").Output()
	if err != nil {
		logger.Errorf("Getting CPU load failed: %s", err)
	}
	cpu_load_averages := strings.Split(string(out), " ")
	if len(cpu_load_averages) > 3 {
		cpu_load_1, err := strconv.Atoi(strings.Split(cpu_load_averages[1], ".")[0])
		if err != nil {
			logger.Errorf("Getting 5 min CPU load failed: %s", err)
		}
		logger.Infof("CPU 5 minutes load average: %d", cpu_load_1)
		return cpu_load_1
	} else {
		logger.Errorf("Getting CPU load failed: %s", err)
		return 0
	}

}
