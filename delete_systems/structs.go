package delete_systems

import (
	"time"
)

type DeleteSystemRequest struct {
	MinionName string `json:"minion_name"`
	Token      string `json:"authentication_token"`
}

type Get_System_Request struct {
	Sessionkey  string `xmlrpc:"sessionKey"`
	System_Name string `xmlrpc:"name"`
}

type Delete_System_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	System_ID  int    `xmlrpc:"sid"`
	Type       string `xmlrpc:"cleanupType"`
}

type ListSystemInfo struct {
	Result []struct {
		Name               string
		Id                 int
		Last_Checkin       time.Time
		Outdated_Pkg_Count int
	}
}

type Delete_System_Return struct {
	Result_ID *int `xml:"params>param>value>i4,omitempty"`
}

type Error_xmlrpc struct {
	Fault Error_Fault `xml:"fault"`
}

type Error_Fault struct {
	Error_Data ErrorData `xml:"value,omitempty"`
}

type ErrorData struct {
	Error_Struct Struct `xml:"struct"`
}
