package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
)

func (s *SPMigration_Email_Body) Send_SPmigration_Email() {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(s.Recipients, "SPMigration Notification - Info", "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		logger.Warningln(err.Error())
	}

	s.Host = hostname
	s.Port = 12345

	//err := r.ParseTemplate("template.html", result)
	template_file := fmt.Sprintf("%s/template_spmigration_info.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, s); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			logger.Warningln(err1.Error())
		}
		logger.Infof("SPMigration Info Email sent. %v", ok)
	} else {
		logger.Warningln(err.Error())
	}

}

func (s *SPMigration_Email_Body) Send_SPmigration_Results() {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(s.Recipients, "SPMigration Notification - Result", "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		logger.Warningln(err.Error())
	}

	s.Host = hostname
	s.Port = 12345

	/* s.SPmigration_Tracking_File = "/var/log/sumapatch/spmigration_t7udp_20230911074532.yaml"
	s.Template_dir = "/srv/jobmonitor/" */
	//err := r.ParseTemplate("template.html", result)
	logger.Infof("Reading JSON file for Email Notification: %s\n", s.SPmigration_Tracking_File)
	targets, err := readJSONFile(s.SPmigration_Tracking_File)
	if err != nil {
		logger.Warningln(err.Error())
	}

	template_file := fmt.Sprintf("%s/template_spmigration_results.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, targets); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			logger.Warningln(err1.Error())
		}
		logger.Infof("SPMigration Results Email sent. %v", ok)
	} else {
		logger.Warningln(err.Error())
	}

}

func readJSONFile(filename string) (*Target_Minions, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %s", err)
	}

	var targets Target_Minions
	if err := json.Unmarshal(content, &targets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %s", err)
	}

	return &targets, nil
}
