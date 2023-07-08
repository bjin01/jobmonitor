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
	Minion_ID                    int
	Minion_Name                  string
	Assign_Vendor_Channel_Job_ID int
	Pkg_Refresh_Job_ID           int
	Outdated_Pkg_Count           int
	Reboot_Job_ID                int
	SP_Migration_Dry_Run_Job_ID  int
	SP_Migration_Job_ID          int
	PostMigration_Reboot_Job_ID  int
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
	return date.Format("2006-01-02T15:04:05")
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

		var response MethodResponse_ActiveSystems_in_Group
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			log.Fatalf("Unmarshal error: %s\n", err)
		}
		fmt.Printf("response of active systems in group: %v\n", response.Params.Param.Value.Array.Data.Values)
		minion_list = append(minion_list)
	}
	return minion_list
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (m *Target_Minions) Get_Minions(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) {
	method := "systemgroup.listSystemsMinimal"
	active_minion_ids := Get_Active_Minions_in_Group(sessionkey, groupsdata)
	fmt.Printf("active_minion_ids: %v\n", active_minion_ids)

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

		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
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
			return
		}

		//fmt.Printf("response: %v\n", response)if len(response.Params.Param.Value.Array.Data.Values) == 1 {

		if len(response.Params.Param.Value.Array.Data.Values) > 0 {

			for _, valueStruct := range response.Params.Param.Value.Array.Data.Values {
				var minion_data Minion_Data
				// Access specific member values by name
				minion_data.Minion_Name = valueStruct.GetMemberValue("name").(string)
				minion_data.Minion_ID = valueStruct.GetMemberValue("id").(int)
				fmt.Printf("name: %s, id: %d\n", minion_data.Minion_Name, minion_data.Minion_ID)
				if Contains(active_minion_ids, minion_data.Minion_ID) {

					m.Minion_List = append(m.Minion_List, minion_data)
				} else {
					fmt.Printf("Minion %s is not active in group %s\n", minion_data.Minion_Name, group)
				}
			}
		}

	}
	m.Show_Minions()
}

func (s *Target_Minions) Show_Minions() {
	for _, minion := range s.Minion_List {
		fmt.Printf("Minion name: %s, Minion ID: %d\n", minion.Minion_Name, minion.Minion_ID)
	}
}

func (t *Target_Minions) Schedule_Pkg_refresh(sessionkey *auth.SumaSessionKey) {
	method := "system.schedulePackageRefresh"
	for _, minion := range t.Minion_List {
		schedule_pkg_refresh_request := Schedule_Pkg_Refresh_Request{
			Sessionkey:         sessionkey.Sessionkey,
			Minion_ID:          minion.Minion_ID,
			EarliestOccourance: Convert_to_ISO8601_DateTime(time.Now()),
		}

		buf, err := gorillaxml.EncodeClientRequest(method, &schedule_pkg_refresh_request)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}

		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Schedule Pkg Refresh API error: %s\n", err)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}

		var response MethodResponse
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			log.Printf("Failed to parse XML-RPC response: %v", err)
			return
		}
		fmt.Printf("Schedule Pkg Refresh response: %v\n", response)
	}
}
