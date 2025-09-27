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
		logger.Infof("Salt %s grains key and or value is not provided. Skipping.\n", s.Salt_no_upgrade_exception_key)
		return nil
	}

	//url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "grains.get"
	s.Arg = []string{s.Salt_no_upgrade_exception_key}

	// Common request execution
	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	salt_request := Salt_Request{
		Client:   "local",
		Tgt:      s.Online_Minions,
		Tgt_type: "list",
		Fun:      s.SaltCmd,
		Arg:      s.Arg,
	}
	response := salt_request.Execute_Command(url, method, s.Token)
	s.Return = response

	var saltResponse No_Upgrade_Grains_return
	if err := json.Unmarshal(response, &saltResponse); err != nil {
		logger.Infoln("Error decoding JSON:", err)
		return nil
	}

	no_upgrade_exception_minions := []string{}
	for _, minion := range saltResponse.Return {
		logger.Debugf("show raw %s exception minion return: %v\n", s.Salt_no_upgrade_exception_key, minion)
		if minionMap, ok := minion.(map[string]interface{}); ok {
			for k, v := range minionMap {
				if reflect.TypeOf(v).String() == "string" {
					logger.Debugf("k is %s %s value is: %s\n", k, s.Salt_no_upgrade_exception_key, v.(string))
					if v.(string) == s.Salt_no_upgrade_exception_value {
						no_upgrade_exception_minions = append(no_upgrade_exception_minions, k)
					}
				}
				if reflect.TypeOf(v).String() == "bool" {
					if strconv.FormatBool(v.(bool)) == s.Salt_no_upgrade_exception_value {
						logger.Infof("k is %s and %s value is bool: %t\n", k, s.Salt_no_upgrade_exception_key, v.(bool))
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
