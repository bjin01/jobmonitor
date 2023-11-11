package pkg_updates

import (
	"time"

	"gorm.io/gorm"
)

type ListMigrationTarget_Request struct {
	Sessionkey                          string `xmlrpc:"sessionKey"`
	Sid                                 int    `xmlrpc:"sid"`
	ExcludeTargetWhereMissingSuccessors bool   `xmlrpc:"excludeTargetWhereMissingSuccessors"`
}

type ListMigrationTarget_Response struct {
	Result []struct {
		Ident    string `xmlrpc:"ident,omitempty"`
		Friendly string `xmlrpc:"friendly,omitempty"`
	}
}

type ListAllChannels_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
}

type ListAllChannels_Response struct {
	Result []struct {
		Id            int    `xmlrpc:"id,omitempty"`
		Name          string `xmlrpc:"name,omitempty"`
		Label         string `xmlrpc:"label,omitempty"`
		Arch_name     string `xmlrpc:"arch_name,omitempty"`
		Provider_name string `xmlrpc:"provider_name,omitempty"`
		Packages      int    `xmlrpc:"packages,omitempty"`
		Systems       int    `xmlrpc:"systems,omitempty"`
	}
}

//DB columns: ID, Group_Name, T7User, Email
type Group struct {
	gorm.Model
	Group_Name string             `json:"group_name" gorm:"primaryKey"`
	T7User     string             `json:"t7user"`
	Email      []Jobchecker_Email `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

//DB columns: ID, Email, GroupID
type Jobchecker_Email struct {
	gorm.Model
	Email   string `json:"email"`
	GroupID uint   `json:"group_id"`
}

//DB columns: ID, Minion_ID, Minion_Name, Minion_Status, Workflow_Step, JobID, JobStatus, Migration_Stage, Migration_Stage_Status, Target_base_channel, Target_Ident, Target_Optional_Channels, Minion_Groups
type Minion_Data struct {
	gorm.Model
	Minion_ID                int                `json:"Minion_ID"`
	Minion_Name              string             `json:"Minion_Name"`
	Minion_Status            string             `json:"Minion_Status"`
	Minion_Remarks           string             `json:"Minion_Remarks"`
	Clm_Stage                string             `json:"Clm_Stage"`
	Workflow_Step            string             `json:"Workflow_Step"`
	JobID                    int                `json:"JobID"`
	JobStatus                string             `json:"JobStatus"`
	Migration_Stage          string             `json:"Migration_Stage"`
	Migration_Stage_Status   string             `json:"Migration_Stage_Status"`
	Target_base_channel      string             `json:"Target_base_channel"`
	Target_Ident             string             `json:"Target_Ident"`
	Target_Optional_Channels []OptionalChannels `json:"Target_Optional_Channels" gorm:"many2many:Minion_Data_OptionalChannels;"`
	Minion_Groups            []Group            `json:"Minion_Groups" gorm:"many2many:Minion_Data_Groups;"`
}

//DB columns: ID, Channel_Label
type OptionalChannels struct {
	gorm.Model
	Channel_Label    string `json:"channel_label" gorm:"primaryKey"`
	Minion_DataRefer uint   `json:"minion_data_refer"`
}

type Update_Groups struct {
	gorm.Model
	Groups          []string `json:"groups"`
	Minions_to_add  []string `json:"minions_to_add"`
	Sqlite_db       string   `json:"sqlite_db"`
	Qualifying_only bool     `json:"qualifying_only"`
	//Delay                           int              `json:"delay"`
	Timeout                         int    `json:"timeout"`
	GatherJobTimeout                int    `json:"gather_job_timeout"`
	Logfile                         string `json:"logfile"`
	Log_Level                       string `json:"log_level"`
	SaltMaster_Address              string `json:"salt_master_address"`
	SaltApi_Port                    int    `json:"salt_api_port"`
	SaltUser                        string `json:"salt_user"`
	SaltPassword                    string `json:"salt_password"`
	Salt_diskspace_grains_key       string `json:"salt_diskspace_grains_key"`
	Salt_diskspace_grains_value     string `json:"salt_diskspace_grains_value"`
	Salt_no_upgrade_exception_key   string `json:"salt_no_upgrade_exception_key"`
	Salt_no_upgrade_exception_value string `json:"salt_no_upgrade_exception_value"`
	Salt_Prep_State                 string `json:"salt_prep_state"`
	Salt_Post_State                 string `json:"salt_post_state"`
	JobcheckerTimeout               int    `json:"jobchecker_timeout"`
	//Reboot_timeout                  int              `json:"reboot_timeout"`
	JobcheckerEmails []string `json:"jobchecker_emails"`
	Patch_Level      string   `json:"patch_level"`
	//Include_Spmigration             bool             `json:"include_spmigration"`
	T7User                  string           `json:"t7user"`
	Token                   string           `json:"authentication_token"`
	Tracking_file_directory string           `json:"tracking_file_directory"`
	Workflow                []map[string]int `json:"workflow"`
	Assigne_channels        []struct {
		Assigne_Channel Assigne_Channel `json:"assign_channel"`
	} `json:"assign_channels"`

	Target_Products []struct {
		Product Target_Product `json:"product"`
	} `json:"products"`
}

type Workflow_Step struct {
	gorm.Model
	Order int    `json:"order"`
	Name  string `json:"name"`
}

type Assigne_Channel struct {
	Current_base_channel string `xmlrpc:"current_base_channel"`
	New_base_prefix      string `xmlrpc:"new_base_prefix"`
}

type Target_Product struct {
	Name                  string            `json:"name"`
	Ident                 string            `json:"ident"`
	Base_Channel          string            `json:"base_channel_label"`
	Clm_Project_Label     string            `json:"clm_project_label"`
	OptionalChildChannels []OptionalChannel `json:"optionalChildChannels"`
}

type OptionalChannel struct {
	Old_Channel string `xmlrpc:"old_channel"`
	New_Channel string `xmlrpc:"new_channel"`
}

type Generic_Job_Response struct {
	JobID int `xmlrpc:"id"`
}

type Get_System_by_Group_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	GroupName  string `xmlrpc:"systemGroupName"`
}

type MethodResponse_ActiveSystems_in_Group struct {
	Params Params_ActiveSystems_in_Group `xml:"params"`
}

type Params_ActiveSystems_in_Group struct {
	Param Param_ActiveSystems_in_Group `xml:"param"`
}

type Param_ActiveSystems_in_Group struct {
	Value Value_ActiveSystems_in_Group `xml:"value"`
}

type Value_ActiveSystems_in_Group struct {
	Array Array_ActiveSystems_in_Group `xml:"array"`
}

type Array_ActiveSystems_in_Group struct {
	Data Data_ActiveSystems_in_Group `xml:"data"`
}

type Data_ActiveSystems_in_Group struct {
	Values []int `xml:"value>i4"`
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

type Job_Chain struct {
	SP_Migration_Hosts []Host_Job_Info
}

type Host_Job_Info struct {
	Assigne_Channels_Job     Assigne_Channels_Job     `json:"Assigne_Channels_Job"`
	Pkg_Refresh_Job          Pkg_Refresh_Job          `json:"Pkg_Refresh_Job"`
	Update_Pkg_Job           Update_Pkg_Job           `json:"Update_Pkg_Job"`
	Reboot_Pre_MigrationJob  Reboot_Pre_MigrationJob  `json:"Reboot_Pre_MigrationJob"`
	SP_Migration_DryRun_Job  SP_Migration_DryRun_Job  `json:"SP_Migration_DryRun_Job"`
	SP_Migration_Job         SP_Migration_Job         `json:"SP_Migration_Job"`
	Reboot_Post_MigrationJob Reboot_Post_MigrationJob `json:"Reboot_Post_MigrationJob"`
	Channel_Environment      string                   `json:"Channel_Environment"`
}

type Assigne_Channels_Job struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type Pkg_Refresh_Job struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type Update_Pkg_Job struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type Reboot_Pre_MigrationJob struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type SP_Migration_DryRun_Job struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type SP_Migration_Job struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type Reboot_Post_MigrationJob struct {
	JobID     int    `json:"JobID"`
	JobStatus string `json:"JobStatus"`
}

type All_Minions_In_Group struct {
	Minion_List map[string][]string
}
