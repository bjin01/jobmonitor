package saltapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *Salt_Data) Login() {

	url := fmt.Sprintf("http://%s:%d/login", s.SaltMaster, s.SaltApi_Port)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s", "eauth": "sharedsecret"}`, s.Username, s.Password))

	transport := &http.Transport{
		Proxy: nil, // This disables proxy settings
	}

	//client := &http.Client{}
	client := &http.Client{
		Transport: transport,
	}
	/* fmt.Printf("url: %s\n", url)
	fmt.Printf("payload: %v\n", payload) */
	req, err := http.NewRequest(method, url, payload)
	//fmt.Printf("req: %v\n", req)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))
	var result Login_Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	//fmt.Println(PrettyPrint(result.Return[0].Token))
	s.Token = result.Return[0].Token
	fmt.Println("Salt login successful")
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
