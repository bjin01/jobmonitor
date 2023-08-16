package saltapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *Salt_Data) Query_Jid() {

	url := fmt.Sprintf("http://%s:%d/jobs/%s", s.SaltMaster, s.SaltApi_Port, s.Jid)
	method := "GET"

	//fmt.Printf("url: %s\n", url)
	payload := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", s.Token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("Error: %v\n", res.Status)
		s.Return = []byte(fmt.Sprintf("Error: %s, maybe the job is not finished or deleted from job cache already.", res.Status))
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	s.Return = prettyJSON.Bytes()
	//fmt.Printf("prettyJSON: %v\n", string(prettyJSON.Bytes()))

	var saltResponse SaltResponse
	if err := json.Unmarshal(body, &saltResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)

		return
	}

	for _, job := range saltResponse.Return {
		fmt.Println("JID:", job.JID)
		fmt.Println("Function:", job.Function)
		if job.Function == "grains.get" && string_array_contains(job.Arguments, "btrfs:for_patching") {
			for hostname, result := range job.Result {
				fmt.Println("grains.get btrfs:for_patching result:")
				fmt.Println("Hostname:", hostname)
				for key, value := range result.(map[string]interface{}) {
					if key == "return" {
						fmt.Printf("Return Value: %s", value.(string))
					}
				}
				fmt.Println()
			}
			continue
		}

		if len(job.Result) == 0 {
			fmt.Println("No result returned")
			continue
		}

		fmt.Println("job returns:", job.Result)
		for hostname, result := range job.Result {
			fmt.Println("Hostname:", hostname)
			parse_interface(result)
			/* fmt.Println("Minion Overall Result:", result.Success) // This will print the raw JSON for each hostname's result

			for key, minion := range result.Return {

				fmt.Println("Key:", key)
				fmt.Println("Name:", minion.Name)
				fmt.Println("Changes:", minion.Changes)
				fmt.Println("Result:", minion.Result)
				fmt.Println("Comment:", minion.Comment)
				fmt.Println("SLS:", minion.Sls)
				fmt.Println("SLS ID:", minion.Sls_Id)
				fmt.Println()
			}*/
		}
	}
	//fmt.Println(string(s.Return))
	return
}

func parse_interface(data interface{}) {
	switch v := data.(type) {
	case string:
		fmt.Printf("%v\n", v)
	case float64:
		fmt.Printf("%v\n", v)
	case bool:
		fmt.Printf("%v\n", v)
	case []interface{}:
		//fmt.Println("is an array:")
		for i, u := range v {
			fmt.Printf("array key %v: ", i)
			parse_interface(u)
		}
	case map[string]interface{}:
		//fmt.Println("is an object:")
		for i, u := range v {
			fmt.Printf("map key %v: ", i)

			parse_interface(u)
		}
	default:
		fmt.Println("unknown type!")
	}
}

func string_array_contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
