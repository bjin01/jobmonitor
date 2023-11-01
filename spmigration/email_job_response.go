package spmigration

import "github.com/bjin01/jobmonitor/email"

func (e *Email_job_info) Send_Job_Response_Email(email_job email.Job_Email_Body) {

	//here we check if for the same jobid an email has been sent already and avoid sending it again
	for _, job := range e.Jobinfo {
		//logger.Infof("--- Jobinfo: %+v\n", job)
		for jobid, server := range job {
			/* logger.Infof("--- Jobid: %d Server: %s\n", jobid, server)
			logger.Info("--- email_job.Job_Response.Job_ID: ", email_job.Job_Response.Job_ID)
			logger.Info("--- email_job.Job_Response.Server_name: ", email_job.Job_Response.Server_name) */

			if email_job.Job_Response.Server_name == server && email_job.Job_Response.Job_ID == jobid {
				return
			}
		}
	}

	//add the jobid and servername to the jobinfo slice for comparison
	e.Jobinfo = append(e.Jobinfo, map[int]string{email_job.Job_Response.Job_ID: email_job.Job_Response.Server_name})
	//logger.Infof("Jobinfo: %+v\n", e.Jobinfo)

	//send the email in go routine
	go email_job.Send_Job_Response()

}
