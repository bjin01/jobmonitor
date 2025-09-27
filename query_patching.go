package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bjin01/jobmonitor/schedules"
)

func readJSONFile_patching(filename string) (*schedules.Jobstatus, error) {
	full_file_path := fmt.Sprintf("/srv/pillar/sumapatch/%s", filename)
	content, err := os.ReadFile(full_file_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %s", err)
	}

	var jobstatus schedules.Jobstatus
	if err := json.Unmarshal(content, &jobstatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
	}

	return &jobstatus, nil
}
