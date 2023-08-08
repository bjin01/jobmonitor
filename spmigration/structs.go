package spmigration

import "time"

type Migration_Groups struct {
	Groups                  []string `json:"groups"`
	Delay                   int      `json:"delay"`
	Timeout                 int      `json:"timeout"`
	GatherJobTimeout        int      `json:"gather_job_timeout"`
	Logfile                 string   `json:"logfile"`
	SaltMaster_Address      string   `json:"salt_master_address"`
	SaltUser                string   `json:"salt_user"`
	SaltPassword            string   `json:"salt_password"`
	JobcheckerTimeout       int      `json:"jobchecker_timeout"`
	JobcheckerEmails        []string `json:"jobchecker_emails"`
	T7User                  string   `json:"t7user"`
	Token                   string   `json:"authentication_token"`
	Tracking_file_directory string   `json:"tracking_file_directory"`
	Assigne_channels        []struct {
		Assigne_Channel Assigne_Channel `json:"assigne_channel"`
	} `json:"assigne_channels"`

	Target_Products []struct {
		Product Target_Product `json:"product"`
	} `json:"products"`
}

type Assigne_Channel struct {
	Current_base_channel string `xmlrpc:"current_base_channel"`
	New_base_prefix      string `xmlrpc:"new_base_prefix"`
}

type Target_Product struct {
	Name         string `json:"name"`
	Ident        string `json:"ident"`
	Base_Channel string `json:"base_channel_label"`
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
	Assigne_Channels_Job     Assigne_Channels_Job
	Pkg_Refresh_Job          Pkg_Refresh_Job
	Update_Pkg_Job           Update_Pkg_Job
	Reboot_Pre_MigrationJob  Reboot_Pre_MigrationJob
	SP_Migration_DryRun_Job  SP_Migration_DryRun_Job
	SP_Migration_Job         SP_Migration_Job
	Reboot_Post_MigrationJob Reboot_Post_MigrationJob
}

type Assigne_Channels_Job struct {
	JobID     int
	JobStatus string
}

type Pkg_Refresh_Job struct {
	JobID     int
	JobStatus string
}

type Update_Pkg_Job struct {
	JobID     int
	JobStatus string
}

type Reboot_Pre_MigrationJob struct {
	JobID     int
	JobStatus string
}

type SP_Migration_DryRun_Job struct {
	JobID     int
	JobStatus string
}

type SP_Migration_Job struct {
	JobID     int
	JobStatus string
}

type Reboot_Post_MigrationJob struct {
	JobID     int
	JobStatus string
}
