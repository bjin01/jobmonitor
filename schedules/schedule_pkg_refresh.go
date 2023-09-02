package schedules

import (
	"log"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Pkg_Refresh_Response struct {
	JobID int `xmlrpc:"id"`
}

func Create_pkg_refresh_job(sessionkey *auth.SumaSessionKey, serverid int, servername string) error {
	method := "system.schedulePackageRefresh"
	schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
		Sessionkey:         sessionkey.Sessionkey,
		Sid:                serverid,
		EarliestOccurrence: time.Now(),
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &schedule_pkg_refresh_request)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", string(buf))
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */
	reply := new(Schedule_Pkg_Refresh_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Printf("Decode Pkg Refresh Job response body failed: %s\n", err)
	}
	log.Printf("%s: Package refresh JobID: %d\n", servername, reply.JobID)
	return nil

}
