package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
)

var Sumahost *string

func MakeRequest(buf []byte) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	api_url := fmt.Sprintf("http://%s/rpc/api", *Sumahost)
	resp, err := client.Post(api_url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
