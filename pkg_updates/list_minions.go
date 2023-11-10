package pkg_updates

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

func Get_Minions(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB) error {
	/* if groupsdata.Target_base_channel == "" {
		logger.Infof("Error: Target base channel is not defined")
		return errors.New("Target base channel is not defined")
	} */

	method := "systemgroup.listSystemsMinimal"
	//active_minion_ids := Get_Active_Minions_in_Group(sessionkey, groupsdata)
	//logger.Infof("active_minion_ids: %v\n", active_minion_ids)
	all_minions := []Minion_Data{}
	//no_target_minions := []Minion_Data{}
	//valid_target_minions := []Minion_Data{}
	for _, group := range groupsdata.Groups {
		all_minions_in_group := make(map[string][]string)
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

			for _, valueStruct := range response.Params.Param.Value.Array.Data.Values {
				var minion_data Minion_Data
				// Access specific member values by name
				minion_group := Group{Group_Name: group}
				minion_data.Minion_Name = valueStruct.GetMemberValue("name").(string)
				minion_data.Minion_ID = valueStruct.GetMemberValue("id").(int)
				minion_data.Minion_Groups = append(minion_data.Minion_Groups, minion_group)
				Delete_Notes(sessionkey, minion_data.Minion_ID)
				result := db.FirstOrCreate(&minion_data, minion_data)
				if result.RowsAffected > 0 {
					logger.Infof("Created minion %s - %d\n", minion_data.Minion_Name, result.RowsAffected)
				} else {
					logger.Infof("Minion %s already exists\n", minion_data.Minion_Name)
				}
				all_minions = append(all_minions, minion_data)
				all_minions_in_group[group] = append(all_minions_in_group[group], minion_data.Minion_Name)

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
			writeMapToYAML(fmt.Sprintf("%s/all_%s_minions.yaml", groupsdata.Tracking_file_directory, group), all_minions_in_group)

		}

	}
	//logger.Infof("all_minions: %v\n", all_minions)

	returned_minions := Detect_Online_Minions(sessionkey, all_minions, groupsdata)

	//logger.Infof("online_minions: %v\n", online_minions)
	for _, minion_data := range returned_minions {
		if minion_data.Minion_ID != 0 && minion_data.Minion_Status == "Online" {
			/* if groupsdata.Include_Spmigration {
				ident, target_migration_base_channel := Find_MigrationTarget_New(sessionkey, minion_data.Minion_ID, groupsdata)
				if ident != "" && target_migration_base_channel != "" {
					minion_data.Target_Ident = ident
					minion_data.Target_base_channel = target_migration_base_channel
					db.Save(&minion_data)
					//valid_target_minions = append(valid_target_minions, minion_data)
					logger.Infof("Minion %s has a valid migration target %s\n",
						minion_data.Minion_Name, target_migration_base_channel)
				} else {
					subject := "no valid migration target"
					body := "minion does not have a valid migration target."
					Add_Note(sessionkey, minion_data.Minion_ID, subject, body)
					db.Save(&minion_data)
					logger.Infof("Minion %s has not a valid migration target\n", minion_data.Minion_Name)
				}
			} else {
				db.Save(&minion_data)

			} */
			db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Status", "Online")
		}

		if minion_data.Minion_ID != 0 && minion_data.Minion_Status == "Offline" {
			minion_data.Minion_Remarks = "Offline"
			db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Status", "Offline")
		}
	}
	return nil
}
