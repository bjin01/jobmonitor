package pkg_updates

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
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

func Update_packages(sessionkey *auth.SumaSessionKey, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	var minion_id_list []int

	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}
		//fmt.Printf("-----------Query DB package update %d\n", result.RowsAffected)
		//logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)
		if stage == Find_Next_Stage(wf, minion) {
			logger.Debugf("Minion %s starts %s stage.\n", minion.Minion_Name, stage)
			minion_id_list = append(minion_id_list, minion.Minion_ID)

		}
	}

	JobID_Pkg_Update := SchedulePackageUpdates(sessionkey, minion_id_list)
	if JobID_Pkg_Update > 0 {
		for _, minion := range minion_list {
			for _, m := range minion_id_list {
				if minion.Minion_ID == m {
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", JobID_Pkg_Update)
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "scheduled")
					db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
					logger.Infof("Minion %s has been scheduled to update packages\n", minion.Minion_Name)
				}
			}
		}
	}
}

func SchedulePackageUpdates(sessionkey *auth.SumaSessionKey, minion_id_list []int) int {
	if len(minion_id_list) == 0 {
		logger.Debugf("No minions to schedule package update\n")
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
		return reply.JobID
	}
	return 0
}
