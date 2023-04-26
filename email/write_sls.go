package email

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/bjin01/jobmonitor/schedules"
)

func Write_SLS(result *schedules.Jobstatus, templates_dir *Templates_Dir) (string, error) {
	template_file := fmt.Sprintf("%s/reboot.sls.template", templates_dir.Dir)
	template, err := template.ParseFiles(template_file)
	current_time := time.Now()
	state_file_dir := "/srv/salt/sumapatch"
	// Capture any error
	if err != nil {
		return "", err
	}

	if result.Reboot_List != "" {
		if _, err := os.Stat(result.Reboot_List); os.IsNotExist(err) {
			return "", err
		}
	}
	if _, err := os.Stat(state_file_dir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Default().Printf("Directory for yaml output file is missing: %s\n", state_file_dir)
		return "", err
	} else {
		fileName_path := fmt.Sprintf("%s/reboot_%s_%s.sls", state_file_dir, result.T7user, current_time.Format("20060102150405"))
		f, err := os.OpenFile(fileName_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		template.Execute(f, result)
		log.Default().Printf("sls file written to %s\n", fileName_path)
		fileName := fmt.Sprintf("reboot_%s_%s", result.T7user, current_time.Format("20060102150405"))
		return fileName, nil
	}
}
