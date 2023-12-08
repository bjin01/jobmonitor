package email

import (
	"fmt"
	"net/smtp"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Job_Email_Body) Send_Pkg_Updates_Email(db *gorm.DB) {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(s.Recipients, "Package Updates Notification - Info", "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		logger.Warningln(err.Error())
	}

	s.Host = hostname
	s.Port = 12345
	minions_list, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
	}

	var groups []Group
	err = db.Preload(clause.Associations).Find(&groups).Error
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	for i := range minions_list {
		minions_list[i].Minion_Groups[0].Ctx_ID = groups[0].Ctx_ID
		//logger.Debugf("Minion in Send_Pkg_Updates_Email: %s - %s\n", minion.Minion_Name, minions_list[i].Minion_Groups[0].Ctx_ID)
	}

	if s.Template_dir == "" {
		s.Template_dir = "/srv/jobmonitor/"
	}

	//err := r.ParseTemplate("template.html", result)
	template_file := fmt.Sprintf("%s/template_pkg_updates_info.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, minions_list); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			logger.Warningln(err1.Error())
		}
		logger.Infof("Package Updates Info Email sent to %v. %v", s.Recipients, ok)
	} else {
		logger.Warningln(err.Error())
	}

}

func (s *Job_Email_Body) Send_Pkg_Updates_Results(db *gorm.DB) {
	auth = smtp.PlainAuth("", "", "", "127.0.0.1")

	r := NewRequest(s.Recipients, "Package Updates Notification - Result", "")
	hostname, err := get_hostname_fqdn()
	if err != nil {
		logger.Warningln(err.Error())
	}

	s.Host = hostname
	s.Port = 12345

	if err != nil {
		logger.Warningln(err.Error())
	}

	minions_list, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
	}

	template_file := fmt.Sprintf("%s/template_pkg_updates_results.html", s.Template_dir)
	if err := r.ParseTemplate(template_file, minions_list); err == nil {
		ok, err1 := r.SendEmail()
		if err1 != nil {
			logger.Warningln(err1.Error())
		}
		logger.Infof("Package Updates Results Email sent. %v", ok)
	} else {
		logger.Warningln(err.Error())
	}

}

func GetAll_Minions_From_DB(db *gorm.DB) ([]Minion_Data, error) {
	var minion_data []Minion_Data
	err := db.Preload(clause.Associations).Find(&minion_data).Error
	//err := db.Model(&grp).Preload("Posts").Find(&grp).Error
	return minion_data, err
}
