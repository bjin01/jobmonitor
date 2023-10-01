package auth

import (
	"github.com/bjin01/jobmonitor/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

type Sumalogin struct {
	Login  string `xmlrpc:"username"`
	Passwd string `xmlrpc:"password"`
}

type SumaSessionKey struct {
	Sessionkey string
}

type SumaLogout struct {
	ReturnInt int
}

func Login(method string, args Sumalogin) (reply SumaSessionKey, err error) {
	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Login API error: %s\n", err)
	}

	err = xml.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		logger.Fatalf("Decode Login response body failed: %s\n", err)
	}
	if resp.StatusCode == 200 && reply.Sessionkey != "" {
		logger.Infoln("SUSE Manager Login successful.")
	}
	//logger.Info("xml decoded %#v\n", reply)
	return
}
