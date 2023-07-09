package groups

import (
	"time"
)

/* type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Minion_ID          int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
} */

type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                uint      `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

/* func (c *Target_Minions) Schedule_Pkg_refresh(sessionkey auth.SumaSessionKey, sumahost string) (err error) {
	api_url := fmt.Sprintf("https://%s/rpc/api", sumahost)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification
		},
	}
	client, _ := xmlrpc.NewClient(api_url, tr)
	defer client.Close()

	//method := "system.scheduleReboot"
	method := "api.getVersion"
	for _, minion := range c.Minion_List {
		//fmt.Println(minion.Minion_Name)
		params := new(Schedule_Pkg_Refresh_Request)
		params.Sessionkey = sessionkey.Sessionkey
		params.Sid = uint(minion.Minion_ID)
		params.EarliestOccurrence = time.Now()

		// Send the XML-RPC request over HTTP
		result := new(string)

		err = client.Call(method, sessionkey.Sessionkey, result)
		if err != nil {
			log.Fatalf("xmlrpc request error: %s\n", err)
		}

		// Print the response
		fmt.Printf("Response: %d\n", result)

	}

	return nil
} */
