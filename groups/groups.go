package groups

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Target_Minions struct {
	Minion_List []Minion_Data
}

type Minion_Data struct {
	Minion_ID              int
	Minion_Name            string
	Host_Job_Info          Host_Job_Info
	Migration_Stage        string
	Migration_Stage_Status string
}

type Schedule_Pkg_Refresh_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

type Schedule_Pkg_Refresh_Response struct {
	JobID int `xmlrpc:"id"`
}

func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	//fmt.Printf("raw time data: %s\n", v)
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

func Convert_to_ISO8601_DateTime(date time.Time) string {
	//convert to ISO8651 date time format

	return date.Format("20060102T15:04:05")
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

func Get_Active_Minions_in_Group(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) []int {
	var minion_list []int
	method := "systemgroup.listActiveSystemsInGroup"
	//method := "systemgroup.listInactiveSystemsInGroup"
	for _, group := range groupsdata.Groups {
		get_system_by_group_request := Get_System_by_Group_Request{
			Sessionkey: sessionkey.Sessionkey,
			GroupName:  group,
		}
		buf, err := gorillaxml.EncodeClientRequest(method, &get_system_by_group_request)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}

		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Get Minions from Group API error: %s\n", err)
		}
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}
		//fmt.Printf("responseBody: %s\n", string(responseBody))
		var response MethodResponse_ActiveSystems_in_Group
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			log.Fatalf("Unmarshal error: %s\n", err)
		}
		//fmt.Printf("response of active systems in group: %v\n", response.Params.Param.Value.Array.Data.Values)
		minion_list = append(minion_list, response.Params.Param.Value.Array.Data.Values...)
	}
	return minion_list
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		//fmt.Printf("acive id: %d - minion_id e: %d\n", a, e)
		if a == e {
			return true
		}
	}
	return false
}

func (m *Target_Minions) Get_Minions(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) error {
	method := "systemgroup.listSystemsMinimal"
	active_minion_ids := Get_Active_Minions_in_Group(sessionkey, groupsdata)
	//fmt.Printf("active_minion_ids: %v\n", active_minion_ids)

	for _, group := range groupsdata.Groups {
		get_system_by_group_request := Get_System_by_Group_Request{
			Sessionkey: sessionkey.Sessionkey,
			GroupName:  group,
		}

		//fmt.Printf("get_system_by_group_request: %v\n", &get_system_by_group_request)
		buf, err := gorillaxml.EncodeClientRequest(method, &get_system_by_group_request)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}

		fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Get Minions from Group API error: %s\n", err)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}
		//fmt.Printf("responseBody: %s\n", responseBody)

		var response MethodResponse
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			log.Printf("Failed to parse XML-RPC response: %v", err)
			return err
		}

		//fmt.Printf("response: %v\n", response)if len(response.Params.Param.Value.Array.Data.Values) == 1 {

		if len(response.Params.Param.Value.Array.Data.Values) > 0 {

			for _, valueStruct := range response.Params.Param.Value.Array.Data.Values {
				var minion_data Minion_Data
				// Access specific member values by name
				minion_data.Minion_Name = valueStruct.GetMemberValue("name").(string)
				minion_data.Minion_ID = valueStruct.GetMemberValue("id").(int)
				//fmt.Printf("name: %s, id: %d\n", minion_data.Minion_Name, minion_data.Minion_ID)
				if Contains(active_minion_ids, minion_data.Minion_ID) {

					m.Minion_List = append(m.Minion_List, minion_data)
				} else {
					fmt.Printf("Minion %s is not active in group %s\n", minion_data.Minion_Name, group)
				}
			}
		}

	}
	//m.Show_Minions()
	return nil
}

func (s *Target_Minions) Show_Minions() {
	for _, minion := range s.Minion_List {
		fmt.Printf("Minion name: %s, Minion ID: %d\n", minion.Minion_Name, minion.Minion_ID)
	}
}

func Orchestrate(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups, sumahost string) {
	var target_minions Target_Minions
	target_minions.Get_Minions(sessionkey, groupsdata)
	//target_minions.Show_Minions()

	/* target_minions.Assign_Channels(sessionkey, groupsdata.Update_Channel_Prefix)
	target_minions.Check_Assigne_Channels_Jobs(sessionkey) // deadline 15min
	target_minions.Schedule_Pkg_refresh(sessionkey)        // pkg refresh
	target_minions.Check_Pkg_Refresh_Jobs(sessionkey)      // deadline 15min */
	JobID_Pkg_Update := target_minions.Schedule_Package_Updates(sessionkey)
	target_minions.Check_Package_Updates_Jobs(sessionkey, JobID_Pkg_Update)
	/* target_minions.Pre_Migration_Reboot()
	target_minions.Check_Pre_Migration_reboot(sessionkey)
	target_minions.SP_Migration_DryRun()
	target_minions.Check_SP_Migration_DryRun()
	target_minions.SP_Migration()
	target_minions.Check_SP_Migration()
	target_minions.Post_Migration_Reboot()
	target_minions.Check_Post_Migration_Reboot_Jobs()
	*/
	/* target_minions.Make_Reports() */

}
