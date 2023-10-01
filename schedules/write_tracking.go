package schedules

import (
	"encoding/json"
	"os"

	"github.com/bjin01/jobmonitor/auth"
)

func Write_Tracking_file(sessionKey *auth.SumaSessionKey, filename string, jobstatus Jobstatus) error {
	// create tracking file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			logger.Fatalf("Error creating tracking file %s: %s\n", filename, err)
		}
		defer file.Close()
	}
	// write tracking file, no append, only write

	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Error opening tracking file %s: %s\n", filename, err)
	}

	err = file.Truncate(0)
	_, err = file.Seek(0, 0)
	// write t struct as json into file
	/* json, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		logger.Fatalf("Error marshalling tracking file: %s\n", err)
	} */
	json, err := json.MarshalIndent(jobstatus, "", "   ")
	if err != nil {
		logger.Fatalf("Error marshalling tracking file %s: %s\n", filename, err)
	}
	if _, err := file.Write(json); err != nil {
		logger.Fatalf("Error writing tracking file %s: %s\n", filename, err)
	}

	defer file.Close()
	return nil
}
