package spmigration

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type Target_Minions struct {
	Minion_List             []Minion_Data       `json:"Minion_List"`
	Tracking_file_name      string              `json:"Tracking_file_name"`
	Suma_Group              string              `json:"Suma_Group"`
	Disk_Check_Disqualified []string            `json:"Disk_Check_Disqualified"`
	No_Upgrade_Exceptions   []string            `json:"No_Upgrade_Exceptions"`
	Offline_Minions         []string            `json:"Offline_Minions"`
	No_Targets_Minions      []Minion_Data       `json:"No_Targets_Minions"`
	CSV_Reports             []string            `json:"CSV_Reports"`
	Jobcheck_Timeout        int                 `json:"Jobcheck_Timeout"`
	Reboot_Timeout          int                 `json:"Reboot_Timeout"`
	Minion_Environment_List []map[string]string `json:"Minion_Environment_List"`
}

type Minion_Data struct {
	Minion_ID                int           `json:"Minion_ID"`
	Minion_Name              string        `json:"Minion_Name"`
	Host_Job_Info            Host_Job_Info `json:"Host_Job_Info"`
	Migration_Stage          string        `json:"Migration_Stage"`
	Migration_Stage_Status   string        `json:"Migration_Stage_Status"`
	Target_base_channel      string        `json:"Target_base_channel"`
	Target_Ident             string        `json:"Target_Ident"`
	Target_Optional_Channels []string      `json:"Target_Optional_Channels"`
}

func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	//logger.Infof("raw time data: %s\n", v)
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
			logger.Fatalf("Encoding error: %s\n", err)
		}

		//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			logger.Fatalf("Get Minions from Group API error: %s\n", err)
		}
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Fatalf("ReadAll error: %s\n", err)
		}
		//logger.Infof("responseBody: %s\n", string(responseBody))
		var response MethodResponse_ActiveSystems_in_Group
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			logger.Fatalf("Unmarshal error: %s\n", err)
		}
		//logger.Infof("response of active systems in group: %v\n", response.Params.Param.Value.Array.Data.Values)
		minion_list = append(minion_list, response.Params.Param.Value.Array.Data.Values...)
	}
	return minion_list
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		//logger.Infof("acive id: %d - minion_id e: %d\n", a, e)
		if a == e {
			return true
		}
	}
	return false
}

func (m *Target_Minions) Get_Minions(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups) error {
	/* if groupsdata.Target_base_channel == "" {
		logger.Infof("Error: Target base channel is not defined")
		return errors.New("Target base channel is not defined")
	} */

	method := "systemgroup.listSystemsMinimal"
	//active_minion_ids := Get_Active_Minions_in_Group(sessionkey, groupsdata)
	//logger.Infof("active_minion_ids: %v\n", active_minion_ids)
	all_minions := []Minion_Data{}
	no_target_minions := []Minion_Data{}
	valid_target_minions := []Minion_Data{}
	for _, group := range groupsdata.Groups {
		get_system_by_group_request := Get_System_by_Group_Request{
			Sessionkey: sessionkey.Sessionkey,
			GroupName:  group,
		}

		//logger.Infof("get_system_by_group_request: %v\n", &get_system_by_group_request)
		buf, err := gorillaxml.EncodeClientRequest(method, &get_system_by_group_request)
		if err != nil {
			logger.Fatalf("Encoding error: %s\n", err)
		}

		//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			logger.Fatalf("Get Minions from Group API error: %s\n", err)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Infof("ReadAll error: %s\n", err)
		}
		//logger.Infof("responseBody: %s\n", responseBody)

		var response MethodResponse
		if err := xml.Unmarshal(responseBody, &response); err != nil {
			logger.Infof("Failed to parse XML-RPC response: %v", err)
			return err
		}

		//logger.Infof("response: %v\n", response)if len(response.Params.Param.Value.Array.Data.Values) == 1 {

		if len(response.Params.Param.Value.Array.Data.Values) > 0 {
			all_minions_in_group := make(map[string][]string)

			for _, valueStruct := range response.Params.Param.Value.Array.Data.Values {
				var minion_data Minion_Data
				// Access specific member values by name
				minion_data.Minion_Name = valueStruct.GetMemberValue("name").(string)
				minion_data.Minion_ID = valueStruct.GetMemberValue("id").(int)
				found_existing_minion := false
				for _, existing_minion := range all_minions {
					if existing_minion.Minion_ID == minion_data.Minion_ID {
						logger.Infof("Minion %s already exists in Minion_List\n", minion_data.Minion_Name)
						found_existing_minion = true
					}
				}
				if !found_existing_minion {
					all_minions = append(all_minions, minion_data)
					Delete_Notes(sessionkey, minion_data.Minion_ID)
					all_minions_in_group[group] = append(all_minions_in_group[group], minion_data.Minion_Name)
				}

			}
			// write all minions in group to a yaml file with group name
			// create tracking file exists
			if _, err := os.Stat(fmt.Sprintf("%s/%s.yaml", groupsdata.Tracking_file_directory, group)); os.IsNotExist(err) {
				file, err := os.Create(fmt.Sprintf("%s/all_%s_minions.yaml", groupsdata.Tracking_file_directory, group))
				if err != nil {
					logger.Infof("Error creating tracking file: %s\n", err)
				}
				defer file.Close()
			}
			// write tracking file, no append, only write
			writeMapToYAML(fmt.Sprintf("%s/all_%s_minions.yaml", groupsdata.Tracking_file_directory, group), all_minions_in_group)

		}

	}
	//logger.Infof("all_minions: %v\n", all_minions)
	online_minions, offline_minions := m.Detect_Online_Minions(sessionkey, all_minions, groupsdata)

	//logger.Infof("online_minions: %v\n", online_minions)
	for _, minion_data := range online_minions {
		if minion_data.Minion_ID != 0 {
			if groupsdata.Include_Spmigration {
				ident, target_migration_base_channel := Find_MigrationTarget(sessionkey, minion_data.Minion_ID, groupsdata)
				if ident != "" && target_migration_base_channel != "" {
					valid_target_minions = append(valid_target_minions, minion_data)
					logger.Infof("Minion %s has a valid migration target %s\n",
						minion_data.Minion_Name, target_migration_base_channel)
				} else {
					no_target_minions = append(no_target_minions, minion_data)
					subject := "no valid migration target"
					body := "minion does not have a valid migration target."
					Add_Note(sessionkey, minion_data.Minion_ID, subject, body)
					logger.Infof("Minion %s has not a valid migration target\n", minion_data.Minion_Name)
				}
			} else {
				valid_target_minions = append(valid_target_minions, minion_data)
			}
		}
	}

	m.Add_Offline_Minions(offline_minions)
	if groupsdata.Include_Spmigration {
		m.Add_Online_Minions(valid_target_minions)
		m.Add_No_Target_Minions(no_target_minions)
	} else {
		m.Add_No_Target_Minions(valid_target_minions)
		logger.Infof("Adding %d online minions for only update workflow.\n", len(m.No_Targets_Minions))
	}
	//m.Show_Minions()
	return nil
}

func (s *Target_Minions) Show_Minions() {
	logger.Infof("Online Minions with targets:\n")
	for _, minion := range s.Minion_List {
		logger.Infof("With Target Minion name: %s, Minion ID: %d\n", minion.Minion_Name, minion.Minion_ID)
	}

	logger.Infof("Online Minions without targets:\n")
	for _, minion := range s.No_Targets_Minions {
		logger.Infof("No Target Minion name: %s, Minion ID: %d\n", minion.Minion_Name, minion.Minion_ID)
	}
}

func Orchestrate(sessionkey *auth.SumaSessionKey, groupsdata *Migration_Groups, sumahost string, email_template_dir string, health *bool) {
	var target_minions Target_Minions
	target_minions.Minion_Environment_List = make([]map[string]string, 0)
	emails := new(email.SPMigration_Email_Body)
	emails.Recipients = groupsdata.JobcheckerEmails

	if groupsdata.JobcheckerTimeout != 0 && groupsdata.JobcheckerTimeout > 50 {
		target_minions.Jobcheck_Timeout = groupsdata.JobcheckerTimeout
		logger.Infof("Set Jobchecker timeout: %d\n", target_minions.Jobcheck_Timeout)
	} else {
		target_minions.Jobcheck_Timeout = 60
		logger.Infof("Use default Jobchecker timeout: %d\n", target_minions.Jobcheck_Timeout)
	}

	if groupsdata.Reboot_timeout != 0 && groupsdata.Reboot_timeout > 20 {
		target_minions.Reboot_Timeout = groupsdata.Reboot_timeout
		logger.Infof("Set Reboot timeout: %d\n", target_minions.Reboot_Timeout)
	} else {
		target_minions.Reboot_Timeout = 50
		logger.Infof("Use default Reboot timeout: %d\n", target_minions.Reboot_Timeout)
	}

	if groupsdata.Tracking_file_directory != "" {
		// create tracking file directory
		if _, err := os.Stat(groupsdata.Tracking_file_directory); os.IsNotExist(err) {
			os.Mkdir(groupsdata.Tracking_file_directory, 0755)
		}
		file_dir := strings.TrimSuffix(groupsdata.Tracking_file_directory, "/")
		if groupsdata.T7User != "" {
			target_minions.Suma_Group = fmt.Sprintf("spmigration_%s_%s", groupsdata.T7User, time.Now().Format("20060102150405"))
			target_minions.Tracking_file_name = fmt.Sprintf("%s/%s.yaml", file_dir, target_minions.Suma_Group)

			emails.T7user = groupsdata.T7User
		} else {
			target_minions.Suma_Group = fmt.Sprintf("spmigration_%s", time.Now().Format("20060102150405"))
			target_minions.Tracking_file_name = fmt.Sprintf("%s/%s.yaml", file_dir, target_minions.Suma_Group)
		}
		// create tracking file

		if _, err := os.Stat(target_minions.Tracking_file_name); os.IsNotExist(err) {
			file, err := os.Create(target_minions.Tracking_file_name)
			if err != nil {
				logger.Fatalf("Error creating tracking file: %s\n", err)
			}

			defer file.Close()
		}
	} else {
		target_minions.Tracking_file_name = "/var/log/spmigration.yaml"
	}
	logger.Infof("Tracking file: %s\n", target_minions.Tracking_file_name)

	target_minions.Get_Minions(sessionkey, groupsdata)
	//logger.Infof("Minions in group: %v\n", target_minions.Minion_List)
	logger.Infof("Qualifying only value set to: %v\n", groupsdata.Qualifying_only)

	//target_minions.Show_Minions()
	target_minions.Write_Tracking_file()
	target_minions.Salt_Refresh_Grains(sessionkey, groupsdata)
	target_minions.Salt_No_Upgrade_Exception_Check(sessionkey, groupsdata)
	target_minions.Salt_Disk_Space_Check(sessionkey, groupsdata)
	target_minions.SPMigration_Group(sessionkey, groupsdata)

	target_minions.Salt_Run_state_apply(sessionkey, groupsdata, "pre")
	target_minions.Show_Minions()
	target_minions.Write_Tracking_file()
	emails.SPmigration_Tracking_File = target_minions.Tracking_file_name
	emails.Template_dir = email_template_dir

	emails.Send_SPmigration_Email()
	if groupsdata.Qualifying_only {
		logger.Infof("Qualifying only is set true so we exit here.\n")
		return
	}
	if groupsdata.Include_Spmigration {
		target_minions.Assign_Channels(sessionkey, groupsdata)
		target_minions.Check_Assigne_Channels_Jobs(sessionkey, health) // deadline 10min
	}

	//target_minions.Schedule_Pkg_refresh(sessionkey)        // pkg refresh
	//target_minions.Check_Pkg_Refresh_Jobs(sessionkey)      // deadline 15min
	JobID_Pkg_Update := target_minions.Schedule_Package_Updates(sessionkey)
	target_minions.Check_Package_Updates_Jobs(sessionkey, JobID_Pkg_Update, health)
	target_minions.Schedule_Reboot(sessionkey)
	target_minions.Check_Reboot_Jobs(sessionkey, health)
	target_minions.Schedule_Pkg_refresh(sessionkey)           // pkg refresh
	target_minions.Check_Pkg_Refresh_Jobs(sessionkey, health) // deadline 15min

	if groupsdata.Include_Spmigration {
		target_minions.ListMigrationTarget(sessionkey, groupsdata)
		target_minions.Schedule_Migration(sessionkey, groupsdata, true)
		target_minions.Check_SP_Migration(sessionkey, true, health)
		target_minions.Schedule_Migration(sessionkey, groupsdata, false)
		target_minions.Check_SP_Migration(sessionkey, false, health)
		target_minions.Salt_Set_Patch_Level(sessionkey, groupsdata)
		target_minions.Salt_Refresh_Grains(sessionkey, groupsdata)
		target_minions.Schedule_Reboot(sessionkey)
		target_minions.Check_Reboot_Jobs(sessionkey, health)
		target_minions.Analyze_Pending_SPMigration(sessionkey, groupsdata, email_template_dir, health)
	} else {
		target_minions.Salt_Set_Patch_Level(sessionkey, groupsdata)
		target_minions.Salt_Refresh_Grains(sessionkey, groupsdata)
	}

	target_minions.Salt_CSV_Report(sessionkey, groupsdata)
	target_minions.Write_Tracking_file()
	target_minions.Salt_Run_state_apply(sessionkey, groupsdata, "post")
	emails.Send_SPmigration_Results()
}

func (t *Target_Minions) Write_Tracking_file() {
	// create tracking file exists
	if _, err := os.Stat(t.Tracking_file_name); os.IsNotExist(err) {
		file, err := os.Create(t.Tracking_file_name)
		if err != nil {
			logger.Fatalf("Error creating tracking file: %s\n", err)
		}
		defer file.Close()
	}
	// write tracking file, no append, only write

	file, err := os.OpenFile(t.Tracking_file_name, os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Error opening tracking file: %s\n", err)
	}

	err = file.Truncate(0)
	_, err = file.Seek(0, 0)
	// write t struct as json into file
	/* json, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		logger.Fatalf("Error marshalling tracking file: %s\n", err)
	} */
	json, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		logger.Fatalf("Error marshalling tracking file: %s\n", err)
	}
	if _, err := file.Write(json); err != nil {
		logger.Fatalf("Error writing tracking file: %s\n", err)
	}

	defer file.Close()
}
