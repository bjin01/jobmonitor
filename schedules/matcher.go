package schedules

import (
	"log"

	"github.com/bjin01/jobmonitor/auth"
)

func (j *Jobstatus) Compare(Sessionkey *auth.SumaSessionKey, scheduled_jobs []Job) {
	listjobs := new(ListJobs)
	err := listjobs.GetPendingjobs(Sessionkey)
	if err != nil {
		log.Default().Printf("GetPendingjobs error: %s\n", err)
	}

	for _, b := range scheduled_jobs {
		for _, y := range listjobs.Pending.Result {
			if y.Id == b.JobID {
				j.Pending = append(j.Pending, b)
			}
		}
	}

	if len(scheduled_jobs) > 0 {
		err = listjobs.GetFailedJobs(Sessionkey)
		if err != nil {
			log.Default().Printf("GetFailedJobs error: %s\n", err)
		}
		for _, b := range scheduled_jobs {
			for _, y := range listjobs.Failed.Result {
				if y.Id == b.JobID {
					j.Failed = append(j.Failed, b)
				}
			}
		}
	}

	if len(scheduled_jobs) > 0 {
		err = listjobs.GetCompletedJobs(Sessionkey)
		if err != nil {
			log.Default().Printf("GetCompletedJobs error: %s\n", err)
		}
		for _, b := range scheduled_jobs {
			for _, y := range listjobs.Completed.Result {
				if y.Id == b.JobID {
					//err = create_pkg_refresh_job(Sessionkey, b.Hostname)
					j.Completed = append(j.Completed, b)
				}
			}
		}
	}

	if len(scheduled_jobs) > 0 {
		for _, b := range scheduled_jobs {

			if !(isExists(b.JobID, j)) {
				log.Printf("append %+v\n", b)
				j.Cancelled = append(j.Cancelled, b)
				continue
			}
		}

	}
}

func isExists(id int, list *Jobstatus) bool {
	for _, l := range list.Pending {
		if l.JobID == id {
			log.Println(l, " Pending")
			return true
		}
	}

	for _, l := range list.Failed {
		if l.JobID == id {
			log.Println(l, " Failed")
			return true
		}
	}

	for _, l := range list.Completed {
		if l.JobID == id {
			log.Println(l, " Completed")
			return true
		}
	}
	return false
}
