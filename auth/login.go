package auth

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
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
	buf, _ := gorillaxml.EncodeClientRequest(method, &args)

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Login API error: %s\n", err)
	}

	defer resp.Body.Close()
	body_bytes, _ := ioutil.ReadAll(resp.Body)

	bodyReader := bytes.NewReader(body_bytes)
	err = gorillaxml.DecodeClientResponse(bodyReader, &reply)
	if err != nil {
		var fault Fault
		err = xml.Unmarshal(body_bytes, &fault)
		if err != nil {
			logger.Printf("Decode Login error response body failed: %s\n", err)
		}
		logger.Printf("Login failed error: %s\n", fault.FaultString)
		return reply, fmt.Errorf("Login to SUMA failed: %s", fault.FaultString)
	}
	if resp.StatusCode == 200 && reply.Sessionkey != "" {
		logger.Infoln("SUSE Manager Login successful.")
	}
	//logger.Info("xml decoded %#v\n", reply)
	return
}
