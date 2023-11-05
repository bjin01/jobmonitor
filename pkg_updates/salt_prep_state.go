package pkg_updates

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/saltapi"
	"gorm.io/gorm"
)

func Salt_Run_state_apply(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, stage string, db *gorm.DB) {
	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.SaltCmd = "state.apply"
	if stage == "pre" {
		saltdata.Arg = []string{groupsdata.Salt_Prep_State}
	} else if stage == "post" {
		saltdata.Arg = []string{groupsdata.Salt_Post_State}
	} else {
		logger.Infof("Salt_Run_state_apply stage is not pre or post. Exiting.\n")
		return
	}

	all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	for _, minion := range all_minions {
		if minion.Minion_Status == "Online" {
			saltdata.Online_Minions = append(saltdata.Online_Minions, minion.Minion_Name)
		}
	}

	url := fmt.Sprintf("http://%s:%d/", saltdata.SaltMaster, saltdata.SaltApi_Port)
	method := "POST"

	if len(saltdata.Online_Minions) > 0 {
		logger.Infof("Salt_Run_state_apply Online_Minions: %s\n", saltdata.Online_Minions)
	} else {
		logger.Infof("Salt_Run_state_apply Online_Minions is empty\n")
		saltdata.Return = []byte("Salt_Run_state_apply Online_Minions is empty")
		return
	}

	saltdata.Login()

	salt_request := saltapi.Salt_Request_Async{
		Tgt:      saltdata.Online_Minions,
		Tgt_type: "list",
		Fun:      saltdata.SaltCmd,
		Arg:      []string{},
	}

	if len(saltdata.Arg) > 0 {
		salt_request.Arg = saltdata.Arg
	} else {
		logger.Infof("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/minions", saltdata.SaltMaster, saltdata.SaltApi_Port)
	response := salt_request.Execute_Command_Async(url, method, saltdata.Token)
	//logger.Infoln(string(response))

	async_response := new(saltapi.Salt_Async_Response)
	if err := json.Unmarshal(response, &async_response); err != nil { // Parse []byte to go struct pointer
		logger.Infof("Can not unmarshal JSON")
	} else {
		logger.Infof("salt state.apply jid: %s\n", async_response.Return[0].Jid)
		saltdata.Jid = async_response.Return[0].Jid
	}
	saltdata.Return = response

	time.Sleep(time.Second * 2)

	deadline := time.Now().Add(time.Duration(15) * time.Minute)

	for time.Now().Before(deadline) {
		err := saltdata.Query_Jid()
		if err != nil {
			logger.Infoln(err)
			logger.Infoln("We will retry salt job query in 5 seconds and final deadline is ", deadline)
		} else {
			deadline = time.Now()
		}
		time.Sleep(time.Second * 5)
	}
}
