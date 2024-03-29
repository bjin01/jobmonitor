package saltapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Btrfs_for_patching_return struct {
	Return []interface{} `json:"return"`
}

func (s *Salt_Data) Run_Disk_Space_Check() []string {
	if s.Salt_diskspace_grains_key == "" || s.Salt_diskspace_grains_value == "" {
		logger.Infof("Salt disk space check is not configured. Skipping.\n")
		return nil
	}

	parent_diskspace_grains_key := strings.Split(s.Salt_diskspace_grains_key, ":")[0]
	s.Delete_Grains_keys(parent_diskspace_grains_key)
	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "grains.get"
	s.Arg = []string{s.Salt_diskspace_grains_key}

	if len(s.Online_Minions) > 0 {
		logger.Infof("Run disk space check for Online_Minions: %s\n", s.Online_Minions)
	} else {
		logger.Debugf("Online_Minions is empty\n")
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
		logger.Debugf("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	logger.Infoln(string(response))
	s.Return = response

	var saltResponse Btrfs_for_patching_return
	if err := json.Unmarshal(response, &saltResponse); err != nil {
		logger.Infoln("Error decoding JSON:", err)
		return nil
	}

	disqualified_minions := []string{}
	for _, minion := range saltResponse.Return {
		if minion.(map[string]interface{}) != nil {
			for k, v := range minion.(map[string]interface{}) {
				if v.(string) != s.Salt_diskspace_grains_value {
					disqualified_minions = append(disqualified_minions, k)
				}
			}
		}
	}

	if len(disqualified_minions) > 0 {
		return disqualified_minions
	} else {
		return nil
	}
}

func (s *Salt_Data) Delete_Grains_keys(grains_key string) {
	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "grains.delkey"
	s.Arg = []string{grains_key, "force=True"}

	if len(s.Online_Minions) > 0 {
		logger.Debugf("Run grains.delkey %s for Online_Minions: %s\n", s.Arg, s.Online_Minions)
	} else {
		logger.Debugf("Online_Minions is empty\n")
		s.Return = []byte("Online_Minions is empty")
		return
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
		logger.Debugf("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	//logger.Infoln(string(response))
	s.Return = response
}

func (s *Salt_Data) Run_Refresh_Grains() {
	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "saltutil.refresh_grains"

	if len(s.Online_Minions) > 0 {
		logger.Infof("Run refresh grains for Online_Minions: %s\n", s.Online_Minions)
	} else {
		logger.Debugf("Online_Minions is empty\n")
		s.Return = []byte("Online_Minions is empty")
		return
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
		logger.Debugf("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	//logger.Infoln(string(response))
	s.Return = response
}
