package saltapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *Salt_Data) Query_Jid() error {
	if s.Jid == "" {
		logger.Infof("Jid is empty\n")
		s.Return = []byte("Jid is empty")
		return nil
	}

	url := fmt.Sprintf("http://%s:%d/jobs/%s", s.SaltMaster, s.SaltApi_Port, s.Jid)
	method := "GET"

	//logger.Infof("url: %s\n", url)
	payload := strings.NewReader("")

	transport := &http.Transport{
		Proxy: nil, // This disables proxy settings
	}

	//client := &http.Client{}
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		logger.Infoln(err)
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", s.Token)

	res, err := client.Do(req)
	if err != nil {
		logger.Infoln(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logger.Infof("Error: %v\n", res.Status)
		s.Return = []byte(fmt.Sprintf("Error: %s, maybe the job is not finished or deleted from job cache already.", res.Status))
		return fmt.Errorf("Error: %s, maybe the job is not finished or deleted from job cache already.", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Infoln(err)
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		logger.Infoln("Error formatting JSON:", err)
		return err
	}

	s.Return = prettyJSON.Bytes()
	//logger.Infof("prettyJSON: %v\n", string(prettyJSON.Bytes()))

	var saltResponse SaltResponse
	if err := json.Unmarshal(body, &saltResponse); err != nil {
		logger.Infoln("Error decoding JSON:", err)

		return err
	}

	for _, job := range saltResponse.Return {
		logger.Infoln("JID:", job.JID)
		logger.Infoln("Function:", job.Function)
		if job.Function == "grains.get" && string_array_contains(job.Arguments, "btrfs:for_patching") {
			for hostname, result := range job.Result {
				logger.Infoln("grains.get btrfs:for_patching result:")
				logger.Infoln("Hostname:", hostname)
				for key, value := range result.(map[string]interface{}) {
					if key == "return" {
						logger.Infof("Return Value: %s", value.(string))
					}
				}
				logger.Infoln()
			}
			continue
		}

		if job.Function == "state.apply" {
			logger.Infoln("State.apply result:")
			logger.Infof("number of targets: %d, number of results %d\n", len(job.Target), len(job.Result))
			if len(job.Target) == len(job.Result) {
				for hostname, result := range job.Result {

					logger.Infoln("Hostname:", hostname)
					value_result := seek_interface_keyval(result, "success", false)
					logger.Infof("value_result: %v\n", value_result)
					//parse_interface(value_result)
				}
				logger.Infoln("All minions returned")
			} else {
				for hostname, result := range job.Result {
					logger.Infoln("Hostname:", hostname)
					value_result := seek_interface_keyval(result, "success", false)
					logger.Infof("value_result: %v\n", value_result)

				}
				logger.Infoln("Still waiting for other minions to return.")
				return fmt.Errorf("Status code: %s, we will retry until all minions returned.", res.Status)
			}
			continue
		}

		if len(job.Result) == 0 {
			logger.Infoln("No result returned")
			continue
		}

		logger.Infoln("job returns:", job.Result)
		for hostname, result := range job.Result {
			logger.Infoln("Hostname:", hostname)
			parse_interface(result)
			/* logger.Infoln("Minion Overall Result:", result.Success) // This will print the raw JSON for each hostname's result

			for key, minion := range result.Return {

				logger.Infoln("Key:", key)
				logger.Infoln("Name:", minion.Name)
				logger.Infoln("Changes:", minion.Changes)
				logger.Infoln("Result:", minion.Result)
				logger.Infoln("Comment:", minion.Comment)
				logger.Infoln("SLS:", minion.Sls)
				logger.Infoln("SLS ID:", minion.Sls_Id)
				logger.Infoln()
			}*/
		}
	}
	//logger.Infoln(string(s.Return))
	return nil
}

func parse_interface(data interface{}) {
	switch v := data.(type) {
	case string:
		logger.Infof("%v\n", v)
	case float64:
		logger.Infof("%v\n", v)
	case bool:
		logger.Infof("%v\n", v)
	case []interface{}:
		//logger.Infoln("is an array:")
		for _, u := range v {
			//logger.Infof("array key %v: ", i)
			parse_interface(u)
		}
	case map[string]interface{}:
		//logger.Infoln("is an object:")
		for _, u := range v {
			//logger.Infof("map key %v: ", i)

			parse_interface(u)
		}
	default:
		logger.Infoln("unknown type!")
	}
}

func seek_interface_keyval(data interface{}, key string, found bool) interface{} {
	switch v := data.(type) {
	case string:

		if found {
			//logger.Infof("return %v\n", v)
			return v
		}
	case float64:
		if found {
			//logger.Infof("return %v\n", v)
			return v
		}
	case bool:
		if found {
			//logger.Infof("return %v\n", v)
			return v
		}
	case []interface{}:
		//logger.Infoln("is an array:")
		for _, u := range v {
			parse_interface(u)

		}
	case map[string]interface{}:
		//logger.Infoln("is an object:")
		for i, u := range v {
			if found {
				//logger.Infof("return %v\n", u)
				return u
			}

			if key == i {
				found = true
				parse_interface(u)
				return u
			}
		}
	default:
		logger.Infoln("unknown type!")
		return nil
	}
	return nil
}

func string_array_contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
