package pkg_updates

import (
	"github.com/bjin01/jobmonitor/email"
	"gorm.io/gorm"
)

func Send_Email(groupsdata *Update_Groups, email_template_dir *email.Templates_Dir, db *gorm.DB) {
	emails := new(email.SPMigration_Email_Body)
	emails.Recipients = groupsdata.JobcheckerEmails
	email_job := new(email.Job_Email_Body)
	email_job.Recipients = groupsdata.JobcheckerEmails
	email_job.Template_dir = email_template_dir.Dir
	email_job.T7user = groupsdata.T7User

	if len(groupsdata.JobcheckerEmails) != 0 {
		emails.Send_Pkg_Updates_Email(db)
	}
}
