package spmigration

import (
	"fmt"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
)

func (m *Target_Minions) Salt_CSV_Report(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {

	if groupsdata.Patch_Level == "" {
		logger.Infof("Patch Level is not provided. Skipping.\n")
		return
	}

	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword

	for _, group := range groupsdata.Groups {
		logger.Infof("Generating Salt_CSV_Report: %s\n", group)
		input_file := fmt.Sprintf("%s/all_%s_minions.yaml", groupsdata.Tracking_file_directory, group)
		csv_file := fmt.Sprintf("/tmp/%s_%s_%s.csv", groupsdata.T7User, group, time.Now().Format("20060102150405"))
		saltdata.Login()
		csv_report_return := saltdata.Run_CSV_Report(input_file, csv_file)
		if len(csv_report_return) > 0 {
			logger.Infof("Minions CSV report finished: %d returned\n", len(csv_report_return))
			m.CSV_Reports = append(m.CSV_Reports, csv_file)
		}
	}
}
