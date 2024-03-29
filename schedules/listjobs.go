package schedules

import (
	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"

	"github.com/divan/gorilla-xmlrpc/xml"
)

func (l *ListJobs) GetCompletedJobs(Sessionkey *auth.SumaSessionKey) error {
	method := "schedule.listAllCompletedActions"

	buf, err := xml.EncodeClientRequest(method, Sessionkey)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("GetCompletedJobs API error: %s\n", err)
	}

	err = xml.DecodeClientResponse(resp.Body, &l.Completed)
	if err != nil {
		logger.Infof("Decode GetCompletedJobs response body failed: %s\n", err)
	}
	logger.Infof("Total %d completed jobs.\n", len(l.Completed.Result))
	return nil
}

func (l *ListJobs) GetFailedJobs(Sessionkey *auth.SumaSessionKey) error {

	method := "schedule.listFailedActions"

	buf, err := xml.EncodeClientRequest(method, Sessionkey)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("GetFailedJobs API error: %s\n", err)
	}

	err = xml.DecodeClientResponse(resp.Body, &l.Failed)
	if err != nil {
		logger.Infof("Decode GetFailedJobs response body failed: %s\n", err)
	}

	logger.Infof("Total %d failed jobs.\n", len(l.Failed.Result))
	return nil
}

func (l *ListJobs) GetPendingjobs(Sessionkey *auth.SumaSessionKey) error {

	method := "schedule.listInProgressActions"

	buf, err := xml.EncodeClientRequest(method, Sessionkey)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("GetPendingjobs API error: %s\n", err)
	}

	err = xml.DecodeClientResponse(resp.Body, &l.Pending)
	if err != nil {
		logger.Infof("Decode GetPendingjobs response body failed: %s\n", err)
	}

	logger.Infof("Total %d pending jobs.\n", len(l.Pending.Result))
	return nil
}
