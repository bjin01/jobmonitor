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

type MethodResponse struct {
	Params Params `xml:"params"`
}

type Params struct {
	Param Param `xml:"param"`
}

type Param struct {
	Value Value `xml:"value"`
}

type Value struct {
	Array Array `xml:"array"`
}

type Array struct {
	Data Data `xml:"data"`
}

type Data struct {
	Values []Struct `xml:"value>struct"`
}

type Struct struct {
	Members []Member `xml:"member"`
}

type Member struct {
	Name  string     `xml:"name"`
	Value InnerValue `xml:"value"`
}

type CustomTime struct {
	time.Time
}

type InnerValue struct {
	StringValue   *string     `xml:"string,omitempty"`
	IntegerValue  *int        `xml:"i4,omitempty"`
	Int           *int        `xml:"int,omitempty"`
	DateTimeValue *CustomTime `xml:"dateTime.iso8601,omitempty"`
	BooleanValue  *bool       `xml:"bool,omitempty"`
}
