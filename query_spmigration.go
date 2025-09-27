package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Minion struct {
	MinionID             int64       `json:"Minion_ID"`
	MinionName           string      `json:"Minion_Name"`
	HostJobInfo          HostJobInfo `json:"Host_Job_Info"`
	MigrationStage       string      `json:"Migration_Stage"`
	MigrationStageStatus string      `json:"Migration_Stage_Status"`
	TargetBaseChannel    string      `json:"Target_base_channel"`
	TargetIdent          string      `json:"Target_Ident"`
}

type HostJobInfo struct {
	AssigneChannelsJob     Job `json:"Assigne_Channels_Job"`
	PkgRefreshJob          Job `json:"Pkg_Refresh_Job"`
	UpdatePkgJob           Job `json:"Update_Pkg_Job"`
	RebootPreMigrationJob  Job `json:"Reboot_Pre_MigrationJob"`
	SPMigrationDryRunJob   Job `json:"SP_Migration_DryRun_Job"`
	SPMigrationJob         Job `json:"SP_Migration_Job"`
	RebootPostMigrationJob Job `json:"Reboot_Post_MigrationJob"`
}

type Job struct {
	JobID     int64  `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type Configuration struct {
	MinionList       []Minion `json:"Minion_List"`
	TrackingFileName string   `json:"Tracking_file_name"`
	SumaGroup        string   `json:"Suma_Group"`
}

func readJSONFile(filename string) (*Configuration, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %s", err)
	}

	var config Configuration
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
	}

	return &config, nil
}
