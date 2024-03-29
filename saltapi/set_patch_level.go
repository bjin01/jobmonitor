package saltapi

import (
	"fmt"
)

type Salt_Set_Patch_Level_return struct {
	Return []interface{} `json:"return"`
}

func (s *Salt_Data) Run_Set_Patch_Level() string {
	if s.Patch_Level == "" {
		logger.Infof("Patch Level is not provided. Skipping.\n")
		return fmt.Sprintln("Patch Level is not provided. Skipping.")
	}

	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	s.SaltCmd = "postpatching.set_patchlevel"
	s.Arg = append(s.Arg, s.Patch_Level)

	if len(s.Online_Minions) > 0 {
		logger.Infof("Run set patch level for Online_Minions: %s\n", s.Online_Minions)
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
		//logger.Debugf("---------set patch level salt_request.Arg: %v\n", salt_request.Arg)
	} else {
		logger.Debugf("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	//logger.Infoln(string(response))
	s.Return = response
	return fmt.Sprintln(string(response))

}
