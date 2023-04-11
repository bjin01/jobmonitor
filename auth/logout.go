package auth

import (
	"log"

	"github.com/bjin01/jobmonitor/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

func Logout(method string, args SumaSessionKey) error {

	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Logout API error: %s\n", err)
	}

	LogoutResult := new(SumaLogout)
	//fmt.Printf("Raw xml: %s\n", fmt.Sprintln(resp))
	err = xml.DecodeClientResponse(resp.Body, LogoutResult)
	if err != nil {
		log.Fatalf("Decode Logout response body failed: %s\n", err)
	}

	if LogoutResult.ReturnInt == 1 {
		log.Println("Logout successful.")
	} else {
		log.Println("Logout failed.")
	}

	return nil
}
