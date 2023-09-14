package spmigration

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Schedule_high_state_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sids               []int     `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
	Test               bool      `xmlrpc:"test"`
}

func (t *Target_Minions) Schedule_high_state(sessionkey *auth.SumaSessionKey) {
	var minion_id_list []int

	for _, minion := range t.Minion_List {
		minion_id_list = append(minion_id_list, minion.Minion_ID)
	}
	JobID_High_State := t.scheduleHighState(sessionkey, minion_id_list)
	if JobID_High_State > 0 {
		for _, minion := range t.Minion_List {
			log.Printf("minion %s is in stage %s with status %s\n", minion.Minion_Name,
				minion.Migration_Stage, minion.Migration_Stage_Status)

			for i, minion := range t.Minion_List {
				for _, minion_id := range minion_id_list {
					if minion.Minion_ID == minion_id {
						var host_info Host_Job_Info
						host_info.Update_Pkg_Job.JobID = JobID_High_State
						host_info.Update_Pkg_Job.JobStatus = "Scheduled"
						t.Minion_List[i].Host_Job_Info = host_info
						t.Minion_List[i].Migration_Stage = "High_State"
						t.Minion_List[i].Migration_Stage_Status = "Scheduled"

					}
				}
			}
		}
	}

}

func (t *Target_Minions) scheduleHighState(sessionkey *auth.SumaSessionKey, minion_id_list []int) int {
	method := "system.scheduleHighstate"
	schedule_high_state_request := Schedule_high_state_Request{
		Sessionkey:         sessionkey.Sessionkey,
		Sids:               minion_id_list,
		EarliestOccurrence: time.Now(),
		Test:               false,
	}
	//fmt.Printf("schedule_high_state_request: %v\n", schedule_high_state_request)
	JobID := scheduleHighState(sessionkey, method, schedule_high_state_request)
	if JobID > 0 {
		log.Printf("High State Job ID: %d\n", JobID)
		log.Printf("Apply High State Job starts at %s\n", schedule_high_state_request.EarliestOccurrence.Format("2006-01-02 15:04:05"))
	} else {
		log.Printf("High State Job ID is 0\n")
	}
	return JobID
}

func scheduleHighState(sessionkey *auth.SumaSessionKey, method string,
	schedule_high_state_request Schedule_high_state_Request) int {
	//fmt.Printf("schedule_high_state_request: %v\n", schedule_high_state_request)
	buf, err := gorillaxml.EncodeClientRequest(method, &schedule_high_state_request)
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
	reply := new(Generic_Job_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Fatalf("Decode high state Job response body failed: %s\n", err)
	}

	if reply.JobID > 0 {

		return reply.JobID
	} else {
		return 0
	}
}
