package email

import (
	"fmt"
	"net/smtp"
)

func (s *Job_Email_Body) Send_Job_Response() {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	subject := fmt.Sprintf("SUMA Job Response - %s - %s", s.Job_Response.Server_name, s.Job_Response.Job_Status)
	r := NewRequest(s.Recipients, subject, "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		logger.Warningln(err.Error())
	}

	s.Host = hostname
	s.Port = 12345

	//err := r.ParseTemplate("template.html", result)
	template_file := fmt.Sprintf("%s/template_job_response_info.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, s.Job_Response); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			logger.Warningln(err1.Error())
		}
		logger.Infof("SUMA Job Response Info Email sent. %v", ok)
	} else {
		logger.Warningln(err.Error())
	}
}
