package email

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/bjin01/jobmonitor/schedules"
)

func Write_SLS(result *schedules.Jobstatus, templates_dir *Templates_Dir) (string, error) {
	template_file := fmt.Sprintf("%s/reboot.sls.template", templates_dir.Dir)
	template, err := template.ParseFiles(template_file)
	// Capture any error
	if err != nil {
		return "", err
	}

	if result.Reboot_List != "" {
		if _, err := os.Stat(result.Reboot_List); os.IsNotExist(err) {
			return "", err
		}
	}
	if _, err := os.Stat("/srv/salt/sumapatch"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Default().Printf("Directory for yaml output file is missing: %s\n", "/srv/salt/sumapatch")
		return "", err
	} else {
		fileName := fmt.Sprintf("/srv/salt/sumapatch/reboot_%s", result.T7user)
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		template.Execute(f, result)
		log.Default().Printf("sls file written to %s\n", fileName)

		return fileName, nil
	}

}
