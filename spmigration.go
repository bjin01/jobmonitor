package main

import (
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/request"
	"github.com/bjin01/jobmonitor/spmigration"
)

func groups_lookup(SUMAConfig *SUMAConfig, groupsdata *spmigration.Migration_Groups, email_template_dir *email.Templates_Dir, health *bool) {
	if health != nil {
		if !*health {
			logger.Infof("Health check failed. Skipping groups lookup.")
			return
		}
	}

	//logger.Info("SP Migration input data %v\n", groupsdata)
	var sumaconf Sumaconf
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		logger.Infof("SUMAKEY is not set. This might cause error for password decryption.")
	}
	for a, b := range SUMAConfig.SUMA {
		sumaconf.Server = a
		b.Password = Decrypt(key, b.Password)
		sumaconf.Password = b.Password
		sumaconf.User = b.User
		if len(b.Email_to) > 0 {
			sumaconf.Email_to = b.Email_to
		}
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		logger.Fatalln(err)
	}
	email_template_directory_string := email_template_dir.Dir
	spmigration.Orchestrate(SessionKey, groupsdata, string(*request.Sumahost), email_template_directory_string, health)
	//logger.Info("target_minions: %v\n", target_minions)
	//logger.Info("sessionkey: %s\n", SessionKey.Sessionkey)
}
