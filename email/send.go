package email

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os/exec"
	"text/template"

	"github.com/bjin01/jobmonitor/schedules"
)

func Sendit(result *schedules.Jobstatus, templates_dir *Templates_Dir) {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(result.JobcheckerEmails, "Jobchecker Notification", "")
	//err := r.ParseTemplate("template.html", result)
	template_file := fmt.Sprintf("%s/template.html", templates_dir.Dir)
	if err := r.ParseTemplate(template_file, result); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			log.Default().Println(err1.Error())
		}
		log.Printf("Email sent. %v", ok)
	} else {
		log.Default().Println(err.Error())
	}

}

func (s *SPMigration_Email_Body) Send_SPmigration_Email() {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(s.Recipients, "SPMigration Notification", "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		log.Default().Println(err.Error())
	}

	s.Host = hostname
	s.Port = 12345

	//err := r.ParseTemplate("template.html", result)
	template_file := fmt.Sprintf("%s/template_spmigration.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, s); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			log.Default().Println(err1.Error())
		}
		log.Printf("Email sent. %v", ok)
	} else {
		log.Default().Println(err.Error())
	}

}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func get_hostname_fqdn() (string, error) {
	cmd := exec.Command("/usr/bin/hostname", "-f")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error when get_hostname_fqdn: %v", err)
	}
	fqdn := out.String()
	fqdn = fqdn[:len(fqdn)-1] // removing EOL

	return fqdn, nil
}

func (r *Request) SendEmail() (bool, error) {
	hostname, err := get_hostname_fqdn()
	if err != nil {
		log.Default().Printf("Failed to get FQDN: %s\n", err)
	}

	fromName := "SUSE Manager"
	fromEmail := "suma@" + hostname
	from := mail.Address{Name: fromName, Address: fromEmail}
	r.from = from.String()
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "127.0.0.1:25"

	/* if err := smtp.SendMail(addr, auth, "suma1@bo2go.home", r.to, msg); err != nil {
		return false, err
	}
	return true, nil */

	c, err := smtp.Dial(addr)

	if err != nil {
		return false, err
	}

	defer c.Close()
	if err = c.Mail(r.from); err != nil {
		return false, err
	}

	for _, addr := range r.to {
		if err = c.Rcpt(addr); err != nil {
			return false, err
		}
	}

	w, err := c.Data()
	if err != nil {
		return false, err
	}
	_, err = w.Write(msg)
	if err != nil {
		return false, err
	}

	err = w.Close()
	if err != nil {
		return false, err
	}

	err = c.Quit()
	// Or alternatively, send with remote service like Amazon SES
	// err = smtp.SendMail(addr, auth, fromEmail, toEmails, bMsg)
	// Handle response from local postfix or remote service
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
