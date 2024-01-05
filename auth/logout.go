package auth

import (
	"github.com/bjin01/jobmonitor/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

func Logout(method string, args SumaSessionKey) error {

	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Logout API error: %s\n", err)
	}

	LogoutResult := new(SumaLogout)
	//logger.Info("Raw xml: %s\n", fmt.Sprintln(resp))
	err = xml.DecodeClientResponse(resp.Body, LogoutResult)
	if err != nil {
		logger.Printf("Decode Logout response body failed: %s\n", err)
	}

	if LogoutResult.ReturnInt == 1 {
		logger.Infoln("Logout successful.")
	} else {
		logger.Infoln("Logout failed.")
	}

	return nil
}
