package pkg_updates

import (
	"context"
	"time"

	"github.com/bjin01/jobmonitor/email"
	"gorm.io/gorm"
)

func Send_Email(ctx context.Context, groupsdata *Update_Groups, email_template_dir *email.Templates_Dir, db *gorm.DB, health *bool, deadline *time.Time) {
	gr := getGoroutineID()
	logger.Infof("Send_Email: Goroutine ID %d\n", gr)

	for time.Now().Before(*deadline) {
		if !*health {
			logger.Infof("Send_Email can't continue due to SUSE Manager health check failed. Please check the logs. continue after 125 seconds.\n")
			time.Sleep(125 * time.Second)
			continue
		}

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Debugf("Send_Email err: %s %s\n", err, groupsdata.Ctx_ID)
			}
			logger.Infof("Send_Email context: finished %s\n", groupsdata.Ctx_ID)
			return
		default:
			logger.Infof("Send_Email: running %s\n", groupsdata.Ctx_ID)
		}

		email_job := new(email.Job_Email_Body)
		email_job.Recipients = groupsdata.JobcheckerEmails
		email_job.Template_dir = email_template_dir.Dir
		email_job.T7user = groupsdata.T7User

		if len(groupsdata.JobcheckerEmails) != 0 {
			email_job.Send_Pkg_Updates_Email(db)
		}

		if groupsdata.Email_Interval != 0 {
			time.Sleep(time.Duration(groupsdata.Email_Interval) * time.Minute)
		} else {
			time.Sleep(10 * time.Minute)
		}
	}
}
