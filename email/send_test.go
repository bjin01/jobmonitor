package email

import (
	"testing"
	"time"

	"github.com/bjin01/jobmonitor/schedules"
)

func TestSendit(t *testing.T) {
	//obj := new(schedules.Jobstatus)
	obj := new(schedules.Jobstatus)
	pendjob := new(schedules.Job)
	//failjob := new(schedules.Job)
	compjob := new(schedules.Job)

	failjob1 := schedules.Job{Hostname: "pxesap01.bo2go.home", JobID: 975, Masterplan: "abc"}
	failjob2 := schedules.Job{Hostname: "pxesap02.bo2go.home", JobID: 976, Masterplan: "syz"}

	cancjob1 := schedules.Job{Hostname: "pxesap01.bo2go.home", JobID: 975, Masterplan: "abc"}
	cancjob2 := schedules.Job{Hostname: "pxesap02.bo2go.home", JobID: 976, Masterplan: "jie"}
	obj.Pending = append(obj.Pending, *pendjob)
	obj.Failed = append(obj.Failed, failjob1, failjob2)
	obj.Completed = append(obj.Completed, *compjob)
	obj.Cancelled = append(obj.Cancelled, cancjob1, cancjob2)
	obj.Offline_minions = append(obj.Offline_minions, "system1", "system2")
	emails := []string{"bo.jin@suseconsulting.ch", "bo.jin@jinbo01.com"}
	a := time.Now().Local().Format(time.RFC822Z)
	//a := time.Now().Local().Format("20060102150405")
	//20060102150405
	obj.JobStartTime = a

	obj.JobcheckerEmails = emails
	obj.T7user = "t7udp"

	templates_dir := &Templates_Dir{Dir: "/srv/jobmonitor"}
	Sendit(obj, templates_dir)
}
