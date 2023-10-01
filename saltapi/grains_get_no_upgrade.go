package saltapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type No_Upgrade_Grains_return struct {
	Return []interface{} `json:"return"`
}

func (s *Salt_Data) Run_No_Upgrade_Grains_Check() []string {
	if s.Salt_no_upgrade_exception_key == "" || s.Salt_no_upgrade_exception_value == "" {
		logger.Infof("Salt no_upgrade grains key and or value is not provided. Skipping.\n")
		return nil
	}

	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "grains.get"
	s.Arg = []string{s.Salt_no_upgrade_exception_key}

	if len(s.Online_Minions) > 0 {
		logger.Infof("Run no_upgrade grains key check for Online_Minions: %s\n", s.Online_Minions)
	} else {
		logger.Infof("Online_Minions is empty\n")
		s.Return = []byte("Online_Minions is empty")
		return nil
	}

	/* salt_request := Salt_Request{
		Client:   s.Salt_Client_Type,
		Tgt:      s.Online_Minions,
		Tgt_type: "list",
		Fun:      s.SaltCmd,
		Arg:      []string{},
	} */

	salt_request := Salt_Request{
		Client:   "local",
		Tgt:      s.Online_Minions,
		Tgt_type: "list",
		Fun:      s.SaltCmd,
		Arg:      []string{},
	}

	if len(s.Arg) > 0 {
		salt_request.Arg = s.Arg
	} else {
		logger.Infof("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	//logger.Infoln(string(response))
	s.Return = response

	var saltResponse No_Upgrade_Grains_return
	if err := json.Unmarshal(response, &saltResponse); err != nil {
		logger.Infoln("Error decoding JSON:", err)
		return nil
	}

	no_upgrade_exception_minions := []string{}
	for _, minion := range saltResponse.Return {
		//logger.Info("show raw no upgrade exception minion return: %v\n", minion)
		if minion.(map[string]interface{}) != nil {
			for k, v := range minion.(map[string]interface{}) {
				//logger.Info("%s no upgrade exception minion is not string: %s\n", k, reflect.TypeOf(v))
				if reflect.TypeOf(v).String() == "string" {
					if v.(string) == s.Salt_no_upgrade_exception_value {
						no_upgrade_exception_minions = append(no_upgrade_exception_minions, k)
					}
				}
				if reflect.TypeOf(v).String() == "bool" {
					if strconv.FormatBool(v.(bool)) == s.Salt_no_upgrade_exception_value {
						no_upgrade_exception_minions = append(no_upgrade_exception_minions, k)
					}
				}
			}
		}
	}

	if len(no_upgrade_exception_minions) > 0 {
		return no_upgrade_exception_minions
	} else {
		return nil
	}
}
