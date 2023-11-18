package email

import (
	"time"

	"gorm.io/gorm"
)

type Templates_Dir struct {
	Dir string
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type SPMigration_Email_Body struct {
	Host                      string
	Port                      int
	T7user                    string
	Template_dir              string
	SPmigration_Tracking_File string
	Recipients                []string
}

type Job_Email_Body struct {
	Host         string
	Port         int
	T7user       string
	Template_dir string
	Recipients   []string
	Job_Response Job_Response
}

type Job_Response struct {
	Server_name  string
	Base_channel string
	Server_id    int
	T7user       string
	Timestamp    time.Time
	//Message      string
	Job_ID     int
	Job_Status string
}

type Target_Minions struct {
	Minion_List             []Minion_Data `json:"Minion_List"`
	Tracking_file_name      string        `json:"Tracking_file_name"`
	Suma_Group              string        `json:"Suma_Group"`
	Disk_Check_Disqualified []string      `json:"Disk_Check_Disqualified"`
	No_Upgrade_Exceptions   []string      `json:"No_Upgrade_Exceptions"`
	Offline_Minions         []string      `json:"Offline_Minions"`
	No_Targets_Minions      []Minion_Data `json:"No_Targets_Minions"`
	CSV_Reports             []string      `json:"CSV_Reports"`
}

type Minion_Data_SPMigration struct {
	Minion_ID              int           `json:"Minion_ID"`
	Minion_Name            string        `json:"Minion_Name"`
	Host_Job_Info          Host_Job_Info `json:"Host_Job_Info"`
	Migration_Stage        string        `json:"Migration_Stage"`
	Migration_Stage_Status string        `json:"Migration_Stage_Status"`
	Target_base_channel    string        `json:"Target_base_channel"`
	Target_Ident           string        `json:"Target_Ident"`
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
	Target_Optional_Channels []OptionalChannels `json:"Target_Optional_Channels" gorm:"foreignKey:Minion_DataRefer"`
	Minion_Groups            []Group            `json:"Minion_Groups" gorm:"many2many:Minion_Data_Groups;"`
}

//DB columns: ID, Group_Name, T7User, Email
type Group struct {
	gorm.Model
	Group_Name string             `json:"group_name"`
	T7User     string             `json:"t7user"`
	Ctx_ID     string             `json:"context_id"`
	Email      []Jobchecker_Email `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

//DB columns: ID, Email, GroupID
type Jobchecker_Email struct {
	gorm.Model
	Email   string `json:"email"`
	GroupID uint   `json:"group_id"`
}

//DB columns: ID, Channel_Label
type OptionalChannels struct {
	gorm.Model
	Channel_Label    string `json:"channel_label"`
	Minion_DataRefer uint   `json:"minion_data_refer"`
}

type Host_Job_Info struct {
	Assigne_Channels_Job     Assigne_Channels_Job     `json:"Assigne_Channels_Job"`
	Pkg_Refresh_Job          Pkg_Refresh_Job          `json:"Pkg_Refresh_Job"`
	Update_Pkg_Job           Update_Pkg_Job           `json:"Update_Pkg_Job"`
	Reboot_Pre_MigrationJob  Reboot_Pre_MigrationJob  `json:"Reboot_Pre_MigrationJob"`
	SP_Migration_DryRun_Job  SP_Migration_DryRun_Job  `json:"SP_Migration_DryRun_Job"`
	SP_Migration_Job         SP_Migration_Job         `json:"SP_Migration_Job"`
	Reboot_Post_MigrationJob Reboot_Post_MigrationJob `json:"Reboot_Post_MigrationJob"`
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
