package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func performHealthCheck(sumaconfig *SUMAConfig) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	sumaserver := new(string)
	for a := range sumaconfig.SUMA {
		*sumaserver = a
	}

	sumaurl := fmt.Sprintf("https://%s/rhn/manager/api/api/systemVersion", *sumaserver)
	//log.Printf("suma url health check: %s\n", sumaurl)

	resp, err := client.Get(sumaurl)
	if err != nil {
		//log.Println("SUMA Health check - API call failed:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code not OK, %d\n", resp.StatusCode)
	}
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("Response: %s\n", string(responseBody)) */
	return nil
}
