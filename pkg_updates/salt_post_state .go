package pkg_updates

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bjin01/jobmonitor/saltapi"
)

func Salt_Run_Post_State(groupsdata *Update_Groups, post_minion_list []string) {
	saltdata := new(saltapi.Salt_Data)
	saltdata.SaltMaster = groupsdata.SaltMaster_Address
	saltdata.SaltApi_Port = groupsdata.SaltApi_Port
	saltdata.Username = groupsdata.SaltUser
	saltdata.Password = groupsdata.SaltPassword
	saltdata.SaltCmd = "state.apply"
	if groupsdata.Salt_Post_State != "" {
		saltdata.Arg = []string{groupsdata.Salt_Post_State}
	} else {
		logger.Infof("Salt_Run_Post_State in config yaml is not provided. Exiting.\n")
		return
	}

	saltdata.Online_Minions = post_minion_list

	method := "POST"

	if len(saltdata.Online_Minions) > 0 {
		logger.Infof("Salt_Run_Post_State Minions: %s\n", saltdata.Online_Minions)
	} else {
		logger.Infof("Salt_Run_Post_State Minions is empty\n")
		saltdata.Return = []byte("Salt_Run_Post_State Minions is empty")
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

	url := fmt.Sprintf("http://%s:%d/minions", saltdata.SaltMaster, saltdata.SaltApi_Port)
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
