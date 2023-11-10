package pkg_updates

import (
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type GetId_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Name       string `xmlrpc:"name"`
}

type GetId_Response struct {
	Result []struct {
		Id                 int       `xmlrpc:"id"`
		Name               string    `xmlrpc:"name"`
		Last_checkin       time.Time `xmlrpc:"last_checkin"`
		Created            time.Time `xmlrpc:"created"`
		Last_boot          time.Time `xmlrpc:"last_boot"`
		Extra_pkg_count    int       `xmlrpc:"extra_pkg_count"`
		Outdated_pkg_count int       `xmlrpc:"outdated_pkg_count"`
	} `xmlrpc:"result"`
}

func Get_SID(sessionkey *auth.SumaSessionKey, system_name string) int {
	method := "system.getId"
	getid_request := GetId_Request{
		Sessionkey: sessionkey.Sessionkey,
		Name:       system_name,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &getid_request)
	if err != nil {
		logger.Infof("Encoding getid_request error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Infof("Encoding getid_request error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", string(buf))
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))

	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s\n", err)
	}
	logger.Infof("responseBody: %+v\n", responseBody) */
	reply := new(GetId_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Infof("Decode getid_response body failed: %s\n", err)
	}

	//logger.Debugf("GetID reponse: %+v\n", reply.Result[0].Id)
	if len(reply.Result) == 0 {
		logger.Infof("No system found with name %s\n", system_name)
		return 0
	}

	if len(reply.Result) > 1 {
		logger.Infof("More than one system found with name %s\n", system_name)
		return 0
	}

	for _, minion := range reply.Result {
		return minion.Id
	}
	return 0
}
