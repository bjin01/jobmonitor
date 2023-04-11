package request

import (
	"bytes"
	"fmt"
	"net/http"
)

var Sumahost *string

func MakeRequest(buf []byte) (*http.Response, error) {
	api_url := fmt.Sprintf("http://%s/rpc/api", *Sumahost)
	resp, err := http.Post(api_url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
