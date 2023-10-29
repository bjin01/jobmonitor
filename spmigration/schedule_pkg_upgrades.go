package spmigration

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Get_Upgradable_Packages_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Sid        int    `xmlrpc:"systemId"`
}

type Get_Upgradable_Packages_Response struct {
	Result []struct {
		Name          string `xmlrpc:"name,omitempty"`
		From_release  string `xmlrpc:"from_release,omitempty"`
		To_epoch      string `xmlrpc:"to_epoch,omitempty"`
		Arch          string `xmlrpc:"arch,omitempty"`
		To_package_id int    `xmlrpc:"to_package_id,omitempty"`
		From_version  string `xmlrpc:"from_version,omitempty"`
		To_version    string `xmlrpc:"to_version,omitempty"`
		From_arch     string `xmlrpc:"from_arch,omitempty"`
		To_arch       string `xmlrpc:"to_arch,omitempty"`
		From_epoch    string `xmlrpc:"from_epoch,omitempty"`
		To_release    string `xmlrpc:"to_release,omitempty"`
	}
}

type SchedulePackageUpdates_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sids               []int     `xmlrpc:"sids"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type SchedulePackageUpdates_Response struct {
	ActionId int
}

func (t *Target_Minions) Schedule_Package_Updates(sessionkey *auth.SumaSessionKey) int {
	var minion_id_list []int

	for _, minion := range t.Minion_List {
		minion_id_list = append(minion_id_list, minion.Minion_ID)
	}

	for _, minion := range t.No_Targets_Minions {
		minion_id_list = append(minion_id_list, minion.Minion_ID)
	}

	JobID_Pkg_Update := t.schedulePackageUpdates(sessionkey, minion_id_list)
	for _, minion := range t.Minion_List {
		logger.Infof("minion %s is in stage %s with status %s\n", minion.Minion_Name,
			minion.Migration_Stage, minion.Migration_Stage_Status)
		logger.Infof("minion %s has job %d with status %s\n", minion.Minion_Name,
			minion.Host_Job_Info.Update_Pkg_Job.JobID, minion.Host_Job_Info.Update_Pkg_Job.JobStatus)
	}
	for _, minion := range t.No_Targets_Minions {
		logger.Infof("minion %s is in stage %s with status %s\n", minion.Minion_Name,
			minion.Migration_Stage, minion.Migration_Stage_Status)
		logger.Infof("minion %s has job %d with status %s\n", minion.Minion_Name,
			minion.Host_Job_Info.Update_Pkg_Job.JobID, minion.Host_Job_Info.Update_Pkg_Job.JobStatus)
	}
	t.Write_Tracking_file()
	if JobID_Pkg_Update > 0 {
		return JobID_Pkg_Update
	} else {
		return 0
	}
}

func (t *Target_Minions) schedulePackageUpdates(sessionkey *auth.SumaSessionKey, minion_id_list []int) int {
	if len(minion_id_list) == 0 {
		logger.Infof("No minions to schedule package update\n")
		return 0
	}

	method := "system.schedulePackageUpdate"
	params := SchedulePackageUpdates_Request{
		Sessionkey: sessionkey.Sessionkey,
		Sids:       minion_id_list,
		//EarliestOccurrence: time.Now().Add(time.Duration(5) * time.Minute),
		EarliestOccurrence: time.Now(),
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
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
	reply := new(Generic_Job_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Fatalf("Decode Pkg Update Job response body failed: %s\n", err)
	}
	if reply.JobID > 0 {
		logger.Infof("Job %d has been scheduled to update packages on %d minions\n", reply.JobID, len(minion_id_list))
		logger.Infof("Package Update Job starts at %s\n", params.EarliestOccurrence.Format("2006-01-02 15:04:05"))
		for i, minion := range t.Minion_List {
			for _, minion_id := range minion_id_list {
				if minion.Minion_ID == minion_id {
					var host_info Host_Job_Info
					host_info.Update_Pkg_Job.JobID = reply.JobID
					host_info.Update_Pkg_Job.JobStatus = "Scheduled"
					t.Minion_List[i].Host_Job_Info = host_info
					t.Minion_List[i].Migration_Stage = "Pkg_Update"
					t.Minion_List[i].Migration_Stage_Status = "Scheduled"

				}
			}
		}

		for i, minion := range t.No_Targets_Minions {
			for _, minion_id := range minion_id_list {
				if minion.Minion_ID == minion_id {
					var host_info Host_Job_Info
					host_info.Update_Pkg_Job.JobID = reply.JobID
					host_info.Update_Pkg_Job.JobStatus = "Scheduled"
					t.No_Targets_Minions[i].Host_Job_Info = host_info
					t.No_Targets_Minions[i].Migration_Stage = "Pkg_Update"
					t.No_Targets_Minions[i].Migration_Stage_Status = "Scheduled"

				}
			}
		}
		return reply.JobID

	}
	return 0
}

func (m *Minion_Data) Get_Upgradable_Packages(sessionkey *auth.SumaSessionKey) {

	method := "system.listLatestUpgradablePackages"
	params := Get_Upgradable_Packages_Request{
		Sessionkey: sessionkey.Sessionkey,
		Sid:        m.Minion_ID,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
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
	reply := new(Get_Upgradable_Packages_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Fatalf("Decode Pkg Update Job response body failed: %s\n", err)
	}
	if len(reply.Result) > 0 {
		logger.Infof("Minion %s has %d packages to upgrade\n", m.Minion_Name, len(reply.Result))
	}
}
