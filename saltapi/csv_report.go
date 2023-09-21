package saltapi

import (
	"fmt"
	"log"
)

type Salt_CSV_Report_return struct {
	Return []interface{} `json:"return"`
}

func (s *Salt_Data) Run_CSV_Report(input_file string, csv_file string) string {
	if input_file == "" {
		log.Printf("input_file not provided. Skipping.\n")
		return fmt.Sprintln("input_file is not provided. Skipping.")
	}

	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	csv_file = fmt.Sprintf("csv_file=%s", csv_file)
	s.SaltCmd = "post_patching.report"
	s.Arg = []string{input_file, csv_file, "presence_check=True"}

	/* if len(s.Online_Minions) > 0 {
		log.Printf("Run csv report for Online_Minions: %s\n", s.Online_Minions)
	} else {
		log.Printf("Online_Minions is empty\n")
		s.Return = []byte("Online_Minions is empty")
		return fmt.Sprintln("Online_Minions is empty")
	} */

	/* salt_request := Salt_Request{
		Client:   s.Salt_Client_Type,
		Tgt:      s.Online_Minions,
		Tgt_type: "list",
		Fun:      s.SaltCmd,
		Arg:      []string{},
	} */

	salt_request := Salt_Request{
		Client:   "runner",
		Tgt:      s.Online_Minions,
		Tgt_type: "list",
		Fun:      s.SaltCmd,
		Arg:      []string{},
	}

	if len(s.Arg) > 0 {
		salt_request.Arg = s.Arg
	} else {
		log.Printf("salt Argument list is empty\n")
	}

	url = fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command(url, method, s.Token)
	fmt.Println(string(response))
	s.Return = response
	return fmt.Sprintln(string(response))

}
