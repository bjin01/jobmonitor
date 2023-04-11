package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	url := "http://127.0.0.1:12345/jobchecker"
	method := "POST"

	payload := strings.NewReader(`{
    "Patching": [
        {
            "saturn": {
                "Patch Job ID is": 758,
                "event send": true,
                "masterplan": "planA"
            }
        },
        {
            "pxesap01.bo2go.home": {
                "Patch Job ID is": 995,
                "event send": true,
                "masterplan": "planA"
            }
        },
        {
            "pxesap02.bo2go.home": {
                "Patch Job ID is": 996,
                "event send": true,
                "masterplan": "planB"
            }
        }
    ],
    "jobchecker_emails": [
        "bo.jin@jinbo01.com",
        "bo.jin@suseconsulting.ch"
    ],
    "jobchecker_timeout": 2,
    "jobstart_delay": 1,
    "offline_minions": [
        "offsystem1",
        "offsystem2",
        "offsystem3"
    ],
    "t7user": "t7udp"

}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
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
	fmt.Println(string(body))
}
