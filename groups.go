package main

import (
	"log"
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/groups"
	"github.com/bjin01/jobmonitor/request"
)

func groups_lookup(SUMAConfig *SUMAConfig, groupsdata *groups.Migration_Groups) {

	//fmt.Printf("SP Migration input data %v\n", groupsdata)
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
	var target_minions = new(groups.Target_Minions)
	target_minions.Get_Minions(SessionKey, groupsdata)
	//fmt.Printf("target_minions: %v\n", target_minions)
	//fmt.Printf("sessionkey: %s\n", SessionKey.Sessionkey)
}
