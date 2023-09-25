package saltapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	StatusCode int
	Message    string
}

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
	apiResponse := string(body)
	if strings.Contains(apiResponse, "<title>401: Unauthorized</title>") {
		// Extract the status code and message from the HTML content
		statusStart := strings.Index(apiResponse, "<title>") + len("<title>")
		statusEnd := strings.Index(apiResponse, "</title>")
		messageStart := strings.Index(apiResponse, "<body>") + len("<body>")
		messageEnd := strings.Index(apiResponse, "</body>")

		statusCode := 401
		message := ""

		if statusStart != -1 && statusEnd != -1 && messageStart != -1 && messageEnd != -1 {
			statusCodeStr := apiResponse[statusStart:statusEnd]
			message = apiResponse[messageStart:messageEnd]

			if parsedStatusCode, err := fmt.Sscanf(statusCodeStr, "%d", &statusCode); err != nil || parsedStatusCode != 1 {
				statusCode = 401 // Default to 401 if parsing fails
			}
		}

		// Create the ErrorResponse struct
		errorResponse := ErrorResponse{
			StatusCode: statusCode,
			Message:    message,
		}

		// Print the parsed struct
		log.Printf("salt api login failed: %+v\n", errorResponse)
		return
	}

	var result Login_Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	//fmt.Println(PrettyPrint(result.Return[0].Token))
	s.Token = result.Return[0].Token
	log.Println("Salt login successful")
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
