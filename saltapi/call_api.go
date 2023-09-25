package saltapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Salt_Data) Run_Manage_Status() {
	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	if len(s.Target_List) > 0 {
		fmt.Printf("Target_List: %s\n", s.Target_List)
	} else {
		fmt.Printf("Target_List is empty\n")
		return
	}

	salt_request := Salt_Request{
		Client:   "runner",
		Tgt:      s.Target_List,
		Tgt_type: "list",
		Fun:      "manage.status",
		Arg:      []string{},
	}

	salt_request.Arg = append(salt_request.Arg, "timeout=7")
	salt_request.Arg = append(salt_request.Arg, "gather_job_timeout=20")

	fmt.Printf("Now calling execute_command\n")
	response := salt_request.Execute_Command(url, method, s.Token)

	minion_status := Runner_Manage_Status_Response{}
	if err := json.Unmarshal(response, &minion_status); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	} else {
		fmt.Printf("minion_status Up: %v\n", minion_status.Return[0].Up)
		fmt.Printf("minion_status Down: %v\n", minion_status.Return[0].Down)
		s.Online_Minions = minion_status.Return[0].Up
		s.Offline_Minions = minion_status.Return[0].Down
	}
	//fmt.Println(string(body))

}

func (s *Salt_Data) Run() {
	s.Run_Manage_Status()

	url := fmt.Sprintf("http://%s:%d/", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	if len(s.Online_Minions) > 0 {
		log.Printf("Online_Minions: %s\n", s.Target_List)
	} else {
		log.Printf("Online_Minions is empty\n")
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

	salt_request := Salt_Request_Async{
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

	url = fmt.Sprintf("http://%s:%d/minions", s.SaltMaster, s.SaltApi_Port)
	response := salt_request.Execute_Command_Async(url, method, s.Token)
	//fmt.Println(string(response))
	s.Return = response
}

func (u *Salt_Request) Execute_Command(url string, method string, token string) []byte {

	payloadBytes, err := json.MarshalIndent(u, "", "   ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}
	//fmt.Printf("payloadBytes: %v\n", string(payloadBytes))

	payload := bytes.NewReader(payloadBytes)

	transport := &http.Transport{
		Proxy: nil, // This disables proxy settings
	}

	//client := &http.Client{}
	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest(method, url, payload)
	//fmt.Printf("req: %v\n", req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", token)
	//fmt.Printf("req: %v\n", req)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Printf("execute command res body: %v\n", res.Body)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}

func (u *Salt_Request_Async) Execute_Command_Async(url string, method string, token string) []byte {

	payloadBytes, err := json.MarshalIndent(u, "", "   ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}
	/* fmt.Printf("payloadBytes: %v\n", string(payloadBytes)) */

	payload := bytes.NewReader(payloadBytes)

	transport := &http.Transport{
		Proxy: nil, // This disables proxy settings
	}

	//client := &http.Client{}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(method, url, payload)
	//fmt.Printf("req: %v\n", req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", token)
	//fmt.Printf("req: %v\n", req)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	/* async_response := new(Salt_Async_Response)
	if err := json.Unmarshal(body, &async_response); err != nil { // Parse []byte to go struct pointer
		log.Println("Can not unmarshal JSON")
	} else {
		log.Printf("salt api async_response: %v\n", async_response)
	} */

	return body
}
