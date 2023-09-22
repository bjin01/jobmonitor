package email

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

type Target_Minions struct {
	Minion_List             []Minion_Data `json:"Minion_List"`
	Tracking_file_name      string        `json:"Tracking_file_name"`
	Suma_Group              string        `json:"Suma_Group"`
	Disk_Check_Disqualified []string      `json:"Disk_Check_Disqualified"`
	No_Upgrade_Exceptions   []string      `json:"No_Upgrade_Exceptions"`
	Offline_Minions         []string      `json:"Offline_Minions"`
	No_Targets_Minions      []string      `json:"No_Targets_Minions"`
	CSV_Reports             []string      `json:"CSV_Reports"`
}

type Minion_Data struct {
	Minion_ID              int           `json:"Minion_ID"`
	Minion_Name            string        `json:"Minion_Name"`
	Host_Job_Info          Host_Job_Info `json:"Host_Job_Info"`
	Migration_Stage        string        `json:"Migration_Stage"`
	Migration_Stage_Status string        `json:"Migration_Stage_Status"`
	Target_base_channel    string        `json:"Target_base_channel"`
	Target_Ident           string        `json:"Target_Ident"`
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
