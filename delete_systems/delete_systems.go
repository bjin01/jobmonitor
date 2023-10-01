package delete_systems

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	logger.Infof("raw time data: %s\n", v)
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
			//logger.Infof("member name: %s\n", name)
			return member.Value.GetFieldValue()
		}
	}
	return nil
}

func (v InnerValue) GetFieldValue() interface{} {
	if v.StringValue != nil {
		//logger.Infof("String value: %s\n", *v.StringValue)
		return *v.StringValue
	} else if v.IntegerValue != nil {
		//logger.Infof("Found integration value. %d\n", v.IntegerValue)
		return *v.IntegerValue
	} else if v.Int != nil {
		//logger.Infof("Found int value. %d\n", *v.Int)
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
		logger.Infof("Failed to parse XML-RPC response: %v", err)
		return
	}
	logger.Warningf("Error ID: %d\n", Error.Fault.Error_Data.Error_Struct.GetMemberValue("faultCode").(int))
	logger.Warningf("Error Message: %s\n", Error.Fault.Error_Data.Error_Struct.GetMemberValue("faultString").(string))

	return
}

func Delete_System(Sessionkey *auth.SumaSessionKey, deletesystemdata *DeleteSystemRequest, emails_to []string) error {

	//var systeminfo ListSystemInfo
	var systemid int
	get_system_obj := Get_System_Request{
		Sessionkey:  Sessionkey.Sessionkey,
		System_Name: deletesystemdata.MinionName,
	}

	//logger.Infof("get_system_obj %s %s\n", get_system_obj.System_Name, get_system_obj.Sessionkey)

	method := "system.getId"
	buf, err := gorillaxml.EncodeClientRequest(method, &get_system_obj)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}

	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Get SystemID API error: %s\n", err)
	}

	var response MethodResponse
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s\n", err)
	}
	if err := xml.Unmarshal(responseBody, &response); err != nil {
		logger.Warningf("Failed to parse XML-RPC response: %v", err)
		return nil
	}

	// Access the first struct value
	if len(response.Params.Param.Value.Array.Data.Values) == 0 {
		logger.Warningf("%s not found in SUMA.", get_system_obj.System_Name)
		subject := "System deleted from SUSE Manager - system not found"
		message := fmt.Sprintf("System %s does not exist in SUSE Manager.\n", deletesystemdata.MinionName)
		email.Send_system_emails(emails_to, subject, message)
		return nil
	}

	if len(response.Params.Param.Value.Array.Data.Values) > 1 {
		for _, b := range response.Params.Param.Value.Array.Data.Values {
			logger.Infof("System name: %s\n", b.GetMemberValue("name"))
		}
		logger.Warningf("We found more than one system with same name. %v\n",
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

		logger.Infoln("Name:", name)
		logger.Infoln("ID:", id)
		logger.Infoln("Last Checkin:", lastCheckin.(*CustomTime).Format("02/01/2006 15:4:5"))
		logger.Infoln("Outdated Package Count:", outdatedPkgCount)
		systemid = id.(int)

	}

	//systemid := 123
	if systemid > 0 {
		del_system_obj := Delete_System_Request{
			Sessionkey: Sessionkey.Sessionkey,
			System_ID:  systemid,
			Type:       "FORCE_DELETE",
		}

		method = "system.deleteSystem"
		buf, err = gorillaxml.EncodeClientRequest(method, &del_system_obj)
		if err != nil {
			logger.Fatalf("Encoding error: %s\n", err)
		}

		/* mybody := []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<methodResponse><params><param><value><i4>1</i4></value></param></params></methodResponse>}`) */
		resp, err = request.MakeRequest(buf)
		if err != nil {
			logger.Fatalf("Delete System API error: %s\n", err)
		}

		responseBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Fatalf("ReadAll error: %s\n", err)
		}

		if strings.Contains(string(responseBody), "XmlRpcFault") {
			//logger.Infoln("do we get here")
			Handle_Xmlrpc_Error(responseBody)
			subject := "System delete from SUSE Manager - failed"
			message := fmt.Sprintf("System %s delete in SUSE Manager. failed: %s\n", deletesystemdata.MinionName, string(responseBody))
			email.Send_system_emails(emails_to, subject, message)
			return nil
		}
	}
	result := new(Delete_System_Return)
	if err := xml.Unmarshal(responseBody, &result); err != nil {
		logger.Infof("Failed to parse XML-RPC response: %v", err)
		return nil
	}

	logger.Infof("Delete system result %d.\n", *result.Result_ID)
	subject := "System deleted from SUSE Manager"
	message := fmt.Sprintf("System %s has been successfully deleted from SUSE Manager.", deletesystemdata.MinionName)
	email.Send_system_emails(emails_to, subject, message)
	return nil
}
