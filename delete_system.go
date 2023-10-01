package main

import (
	"log"
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/delete_systems"
	"github.com/bjin01/jobmonitor/request"
)

func Delete_System(SUMAConfig *SUMAConfig, deleteSystemdata *delete_systems.DeleteSystemRequest) {
	logger.Infof("deleteSystemdata %s\n", deleteSystemdata.MinionName)
	var sumaconf Sumaconf
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
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
		log.Fatal(err)
	}
	logger.Infof("sessionkey: %s\n", SessionKey.Sessionkey)
	logger.Infof("Deleting System in SUMA: %s\n", deleteSystemdata.MinionName)
	err = delete_systems.Delete_System(SessionKey, deleteSystemdata, sumaconf.Email_to)
	if err != nil {
		log.Fatal(err)
	}

}

func isValidAuthToken(token string) bool {
	logger.Infof("token: %s\n", token)
	if token == os.Getenv("SUMAKEY") {
		return true
	} else {
		return false
	}

}
