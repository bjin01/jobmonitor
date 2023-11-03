package pkg_updates

import (
	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Add_Note_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Sid        int    `xmlrpc:"sid"`
	Subject    string `xmlrpc:"subject"`
	Body       string `xmlrpc:"body"`
}

type Delete_Notes_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Sid        int    `xmlrpc:"sid"`
}

func Add_Note(sessionkey *auth.SumaSessionKey, sid int, subject string, note string) error {
	method := "system.addNote"
	add_note_object := Add_Note_Request{
		Sessionkey: sessionkey.Sessionkey,
		Sid:        sid,
		Subject:    subject,
		Body:       note,
	}
	buf, err := gorillaxml.EncodeClientRequest(method, &add_note_object)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	if resp.StatusCode != 200 {
		logger.Infof("Add_Note error: %s\n", err)
	}
	return nil
}

func Delete_Notes(sessionkey *auth.SumaSessionKey, sid int) error {
	method := "system.deleteNotes"
	delete_nodes_object := Delete_Notes_Request{
		Sessionkey: sessionkey.Sessionkey,
		Sid:        sid,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &delete_nodes_object)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	if resp.StatusCode != 200 {
		logger.Infof("Delete_Notes error: %s\n", err)
	}
	return nil
}
