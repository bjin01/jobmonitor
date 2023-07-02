package delete_systems

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

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

func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	fmt.Printf("raw time data: %s\n", v)
	year, _ := strconv.Atoi(v[0:4])
	month, _ := strconv.Atoi(v[4:6])
	day, _ := strconv.Atoi(v[6:8])
	hour, _ := strconv.Atoi(v[9:11])
	minute, _ := strconv.Atoi(v[12:14])
	second, _ := strconv.Atoi(v[15:17])

	temp_time := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	*c = CustomTime{temp_time}
	return nil
}

func (s Struct) GetMemberValue(name string) interface{} {
	for _, member := range s.Members {
		if member.Name == name {
			//fmt.Printf("member name: %s\n", name)
			return member.Value.GetFieldValue()
		}
	}
	return nil
}

func (v InnerValue) GetFieldValue() interface{} {
	if v.StringValue != nil {
		//fmt.Printf("String value: %s\n", *v.StringValue)
		return *v.StringValue
	} else if v.IntegerValue != nil {
		//fmt.Printf("Found integration value. %d\n", v.IntegerValue)
		return *v.IntegerValue
	} else if v.Int != nil {
		//fmt.Printf("Found int value. %d\n", *v.Int)
		return *v.Int
	} else if v.DateTimeValue != nil {
		return v.DateTimeValue
	} else if v.BooleanValue != nil {
		return *v.BooleanValue
	}
	return nil
}

func Handle_Xmlrpc_Error(message []byte) {
	Error := new(Error_xmlrpc)
	if err := xml.Unmarshal(message, &Error); err != nil {
		fmt.Printf("Failed to parse XML-RPC response: %v", err)
		return
	}
	log.Printf("Error ID: %d\n", Error.Fault.Error_Data.Error_Struct.GetMemberValue("faultCode").(int))
	log.Printf("Error Message: %s\n", Error.Fault.Error_Data.Error_Struct.GetMemberValue("faultString").(string))
	return
}

func Delete_System(Sessionkey *auth.SumaSessionKey, deletesystemdata *DeleteSystemRequest) error {

	//var systeminfo ListSystemInfo
	get_system_obj := Get_System_Request{
		Sessionkey:  Sessionkey.Sessionkey,
		System_Name: deletesystemdata.MinionName,
	}

	//fmt.Printf("get_system_obj %s %s\n", get_system_obj.System_Name, get_system_obj.Sessionkey)

	method := "system.getId"
	buf, err := gorillaxml.EncodeClientRequest(method, &get_system_obj)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}

	fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Get SystemID API error: %s\n", err)
	}

	var response MethodResponse
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	if err := xml.Unmarshal(responseBody, &response); err != nil {
		log.Printf("Failed to parse XML-RPC response: %v", err)
		return nil
	}

	// Access the first struct value
	if len(response.Params.Param.Value.Array.Data.Values) == 0 {
		log.Printf("%s not found in SUMA.", get_system_obj.System_Name)
		//return nil
	}

	if len(response.Params.Param.Value.Array.Data.Values) > 1 {
		for _, b := range response.Params.Param.Value.Array.Data.Values {
			fmt.Printf("System name: %s\n", b.GetMemberValue("name"))
		}
		log.Printf("We found more than one system with same name. %v\n",
			response.Params.Param.Value.Array.Data.Values)
		return nil
	}

	if len(response.Params.Param.Value.Array.Data.Values) == 1 {
		valueStruct := response.Params.Param.Value.Array.Data.Values[0]

		// Access specific member values by name
		name := valueStruct.GetMemberValue("name")
		id := valueStruct.GetMemberValue("id")
		lastCheckin := valueStruct.GetMemberValue("last_checkin")
		outdatedPkgCount := valueStruct.GetMemberValue("outdated_pkg_count")

		log.Println("Name:", name)
		log.Println("ID:", id)
		log.Println("Last Checkin:", lastCheckin.(*CustomTime).Format("02/01/2006 15:4:5"))
		log.Println("Outdated Package Count:", outdatedPkgCount)
	}

	//systemid := id.(int)
	systemid := 123
	if systemid > 0 {
		del_system_obj := Delete_System_Request{
			Sessionkey: Sessionkey.Sessionkey,
			System_ID:  systemid,
			Type:       "FORCE_DELETE",
		}

		method = "system.deleteSystem"
		buf, err = gorillaxml.EncodeClientRequest(method, &del_system_obj)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}

		/* mybody := []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<methodResponse><params><param><value><i4>1</i4></value></param></params></methodResponse>}`) */
		resp, err = request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Delete System API error: %s\n", err)
		}

		responseBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}

		// Print the response as a string
		fmt.Println(string(responseBody))

		if strings.Contains(string(responseBody), "XmlRpcFault") {
			//fmt.Println("do we get here")
			Handle_Xmlrpc_Error(responseBody)
			return nil
		}
	}
	result := new(Delete_System_Return)
	if err := xml.Unmarshal(responseBody, &result); err != nil {
		fmt.Printf("Failed to parse XML-RPC response: %v", err)
		return nil
	}

	/* err = gorillaxml.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		log.Printf("Decode system delete response body failed: %s\n", err)
	} */
	//fmt.Printf("Delete system result %d.\n", *result.Params.Param.Value.Id)
	fmt.Printf("Delete system result %d.\n", *result.Result_ID)
	return nil
}
