package saltapi

import (
	"fmt"
)

type Salt_Refresh_Grains_return struct {
	Return []interface{} `json:"return"`
}

func (s *Salt_Data) Saltutil_Refresh_Grains() string {

	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "saltutil.refresh_grains"

	if len(s.Online_Minions) > 0 {
		logger.Infof("Run saltutil.refresh_grains for Online_Minions: %s\n", s.Online_Minions)
	} else {
		logger.Debugf("Online_Minions is empty\n")
		s.Return = []byte("Online_Minions is empty")
		return fmt.Sprintln("Online_Minions is empty")
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
	return fmt.Sprintln(string(response))

}
