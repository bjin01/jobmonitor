package spmigration

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Target_Minions struct {
	Minion_List        []Minion_Data
	Tracking_file_name string
	Suma_Group         string
}

type Minion_Data struct {
	Minion_ID              int
	Minion_Name            string
	Host_Job_Info          Host_Job_Info
	Migration_Stage        string
	Migration_Stage_Status string
	Target_base_channel    string
	Target_Ident           string
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
	/* if groupsdata.Target_base_channel == "" {
		log.Printf("Error: Target base channel is not defined")
		return errors.New("Target base channel is not defined")
	} */

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

		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Get Minions from Group API error: %s\n", err)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ReadAll error: %s\n", err)
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
					ident, target_migration_base_channel := Find_MigrationTarget(sessionkey, minion_data.Minion_ID, groupsdata)
					if ident != "" && target_migration_base_channel != "" {
						m.Minion_List = append(m.Minion_List, minion_data)
						log.Printf("Minion %s has a valid migration target %s\n",
							minion_data.Minion_Name, target_migration_base_channel)
					} else {
						log.Printf("Minion %s has not a valid migration target\n", minion_data.Minion_Name)
					}
				} else {
					log.Printf("%s is not active in group %s\n", minion_data.Minion_Name, group)
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

func Orchestrate(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups, sumahost string, health *bool) {
	var target_minions Target_Minions
	if groupsdata.Tracking_file_directory != "" {
		// create tracking file directory
		if _, err := os.Stat(groupsdata.Tracking_file_directory); os.IsNotExist(err) {
			os.Mkdir(groupsdata.Tracking_file_directory, 0755)
		}
		file_dir := strings.TrimSuffix(groupsdata.Tracking_file_directory, "/")
		if groupsdata.T7User != "" {
			target_minions.Tracking_file_name = fmt.Sprintf("%s/spmigration_%s_%s.yaml",
				file_dir, groupsdata.T7User, time.Now().Format("20060102150405"))
		} else {
			target_minions.Tracking_file_name = fmt.Sprintf("%s/spmigration_%s.yaml", file_dir,
				time.Now().Format("20060102150405"))
		}
		// create tracking file

		if _, err := os.Stat(target_minions.Tracking_file_name); os.IsNotExist(err) {
			file, err := os.Create(target_minions.Tracking_file_name)
			if err != nil {
				log.Fatalf("Error creating tracking file: %s\n", err)
			}

			defer file.Close()
		}
	} else {
		target_minions.Tracking_file_name = "/var/log/spmigration.yaml"
	}
	log.Printf("Tracking file: %s\n", target_minions.Tracking_file_name)

	target_minions.Get_Minions(sessionkey, groupsdata)
	target_minions.SPMigration_Group(sessionkey, groupsdata)
	//target_minions.Show_Minions()
	target_minions.Write_Tracking_file()
	target_minions.Assign_Channels(sessionkey, groupsdata)
	target_minions.Check_Assigne_Channels_Jobs(sessionkey, health) // deadline 15min
	//target_minions.Schedule_Pkg_refresh(sessionkey)        // pkg refresh
	//target_minions.Check_Pkg_Refresh_Jobs(sessionkey)      // deadline 15min
	JobID_Pkg_Update := target_minions.Schedule_Package_Updates(sessionkey)
	target_minions.Check_Package_Updates_Jobs(sessionkey, JobID_Pkg_Update, health)
	target_minions.Schedule_Reboot(sessionkey)
	target_minions.Check_Reboot_Jobs(sessionkey, health)
	target_minions.Schedule_Pkg_refresh(sessionkey)           // pkg refresh
	target_minions.Check_Pkg_Refresh_Jobs(sessionkey, health) // deadline 15min
	target_minions.ListMigrationTarget(sessionkey, groupsdata)
	target_minions.Schedule_Migration(sessionkey, groupsdata, true)
	target_minions.Check_SP_Migration(sessionkey, true, health)
	target_minions.Schedule_Migration(sessionkey, groupsdata, false)
	target_minions.Check_SP_Migration(sessionkey, false, health)
	target_minions.Schedule_Reboot(sessionkey)
	target_minions.Check_Reboot_Jobs(sessionkey, health)
	target_minions.Analyze_Pending_SPMigration(sessionkey, groupsdata, health)

	/* target_minions.Make_Reports() */

}

func (t *Target_Minions) Write_Tracking_file() {
	// create tracking file exists
	if _, err := os.Stat(t.Tracking_file_name); os.IsNotExist(err) {
		file, err := os.Create(t.Tracking_file_name)
		if err != nil {
			log.Fatalf("Error creating tracking file: %s\n", err)
		}
		defer file.Close()
	}
	// write tracking file, no append, only write

	file, err := os.OpenFile(t.Tracking_file_name, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening tracking file: %s\n", err)
	}

	// write t struct as json into file
	json, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		log.Fatalf("Error marshalling tracking file: %s\n", err)
	}
	if _, err := file.Write(json); err != nil {
		log.Fatalf("Error writing tracking file: %s\n", err)
	}

	defer file.Close()
}
