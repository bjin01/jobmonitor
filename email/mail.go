package email

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/bjin01/jobmonitor/schedules"
)

var host = "127.0.0.1"
var port = "25"
var username = ""
var password = ""
var auth = smtp.PlainAuth("", username, password, host)
var addr = host + ":" + port

func Sendmail(result *schedules.Jobstatus, emailaddr []string) {
	fromName := "SUSE Manager"
	fromEmail := "suma1@bo2go.home"
	toNames := emailaddr
	toEmails := emailaddr
	subject := "This is the subject of your email"
	body := ""
	body += fmt.Sprintf("Pending Jobs:\r\n")
	for _, b := range result.Pending {
		body += fmt.Sprintf("\t%s: %d\r\n", b.Hostname, b.JobID)
	}

	body += fmt.Sprintf("Failed Jobs:\r\n")
	for _, b := range result.Failed {
		body += fmt.Sprintf("\t%s: %d\r\n", b.Hostname, b.JobID)
	}

	body += fmt.Sprintf("Completed Jobs:\r\n")
	for _, b := range result.Completed {
		body += fmt.Sprintf("\t%s: %d\r\n", b.Hostname, b.JobID)
	}

	body += fmt.Sprintf("Cancelled Jobs:\r\n")
	for _, b := range result.Cancelled {
		body += fmt.Sprintf("\t%s: %d\r\n", b.Hostname, b.JobID)
	}

	// Build RFC-2822 email
	toAddresses := []string{}

	for i := range toEmails {
		to := mail.Address{Name: toNames[i], Address: toEmails[i]}
		toAddresses = append(toAddresses, to.String())
	}

	toHeader := strings.Join(toAddresses, ", ")
	from := mail.Address{Name: fromName, Address: fromEmail}
	fromHeader := from.String()
	subjectHeader := subject
	header := make(map[string]string)
	header["To"] = toHeader
	header["From"] = fromHeader
	header["Subject"] = subjectHeader
	header["Content-Type"] = `text/html; charset="UTF-8"`
	msg := ""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg += "\r\n" + body
	bMsg := []byte(msg)
	// Send using local postfix service
	c, err := smtp.Dial(addr)

	if err != nil {
		return
	}

	defer c.Close()
	if err = c.Mail(fromHeader); err != nil {
		return
	}

	for _, addr := range toEmails {
		if err = c.Rcpt(addr); err != nil {
			return
		}
	}

	w, err := c.Data()
	if err != nil {
		return
	}
	_, err = w.Write(bMsg)
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	err = c.Quit()
	// Or alternatively, send with remote service like Amazon SES
	// err = smtp.SendMail(addr, auth, fromEmail, toEmails, bMsg)
	// Handle response from local postfix or remote service
	if err != nil {
		return
	}
}
