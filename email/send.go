package email

import (
	"bytes"
	"log"
	"net/mail"
	"net/smtp"
	"text/template"

	"github.com/bjin01/jobmonitor/schedules"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func Sendit(result *schedules.Jobstatus) {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(result.JobcheckerEmails, "Jobchecker Notification", "")
	//err := r.ParseTemplate("template.html", result)
	if err := r.ParseTemplate("email/template.html", result); err == nil {
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

func (r *Request) SendEmail() (bool, error) {
	fromName := "SUSE Manager"
	fromEmail := "suma1@bo2go.home"
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
