package spmigration

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ScheduleSPMigrationDryRun_Request struct {
	Sessionkey            string    `xmlrpc:"sessionKey"`
	Sid                   int       `xmlrpc:"sid"`
	TargetIdent           string    `xmlrpc:"targetIdent"`
	BaseChannelLabel      string    `xmlrpc:"baseChannelLabel"`
	OptionalChildChannels []string  `xmlrpc:"optionalChildChannels"`
	AllowVendorChange     bool      `xmlrpc:"allowVendorChange"`
	DryRun                bool      `xmlrpc:"dryRun"`
	EarliestOccurrence    time.Time `xmlrpc:"earliestOccurrence"`
}

type ScheduleSPMigrationDryRun_Response struct {
	JobID int `xmlrpc:"id"`
}

func (t *Target_Minions) Schedule_Migration_DryRun(sessionkey *auth.SumaSessionKey, UserData *Migration_Groups) {
	method := "system.scheduleProductMigration"
	for i, minion := range t.Minion_List {
		if minion.Target_Ident == "" {
			log.Default().Printf("Target Ident is empty for minion %s\n", minion.Minion_Name)
			continue
		}
		schedule_spmigration_request := ScheduleSPMigrationDryRun_Request{}
		schedule_spmigration_request.Sessionkey = sessionkey.Sessionkey
		schedule_spmigration_request.Sid = minion.Minion_ID
		schedule_spmigration_request.TargetIdent = minion.Target_Ident
		schedule_spmigration_request.BaseChannelLabel = minion.Target_base_channel
		schedule_spmigration_request.DryRun = true
		schedule_spmigration_request.AllowVendorChange = true
		schedule_spmigration_request.EarliestOccurrence = time.Now()
		buf, err := gorillaxml.EncodeClientRequest(method, &schedule_spmigration_request)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Printf("Encoding scheduleProductMigration error: %s\n", err)
		}
		reply := new(ScheduleSPMigrationDryRun_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			log.Printf("Decode scheduleProductMigration_DryRun Job response body failed: %s\n", err)
		}
		//log.Printf("scheduleProductMigration_DryRun JobID: %d\n", reply.JobID)
		var host_info Host_Job_Info
		host_info.Pkg_Refresh_Job.JobID = reply.JobID
		host_info.Pkg_Refresh_Job.JobStatus = "Scheduled"

		if reply.JobID > 0 {
			t.Minion_List[i].Host_Job_Info = host_info
			t.Minion_List[i].Migration_Stage = "Product Migration DryRun"
			t.Minion_List[i].Migration_Stage_Status = "Scheduled"
			log.Printf("Product migration dryrun JobID: %d %s\n", reply.JobID, minion.Minion_Name)
		} else {
			log.Printf("Minion %s product migration dryrun not possible.\n", minion.Minion_Name)
			continue
		}
	}
}
