package email

import (
	"log"
	"net/smtp"
)

func Send_system_emails(recipients []string, subject string, message string) {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(recipients, subject, message)
	ok, err1 := r.SendEmail()
	if err1 != nil {
		log.Default().Println(err1.Error())
	}
	log.Printf("Email sent. %v", ok)
}
