package schedules

import "time"

type ListCompletedJobs struct {
	Result []struct {
		Scheduler         string
		Name              string
		CompletedSystems  int
		FailedSystems     int
		InProgressSystems int
		Id                int
		Type              string
		Earliest          time.Time
	}
}

type ListFailedJobs struct {
	Result []struct {
		Scheduler         string
		Name              string
		CompletedSystems  int
		FailedSystems     int
		InProgressSystems int
		Id                int
		Type              string
		Earliest          time.Time
	}
}

type ListPendingJobs struct {
	Result []struct {
		Scheduler         string
		Name              string
		CompletedSystems  int
		FailedSystems     int
		InProgressSystems int
		Id                int
		Type              string
		Earliest          time.Time
	}
}

type ListJobs struct {
	Completed          ListCompletedJobs
	Failed             ListFailedJobs
	Pending            ListPendingJobs
	Found_Pending_Jobs bool
}

type Job struct {
	Hostname   string `json:"hostname,omitempty"`
	ServerID   int    `json:"serverid,omitempty"`
	JobID      int    `json:"jobid,omitempty"`
	Masterplan string `json:"masterplan,omitempty"`
}

type Full_Update_Jobs struct {
	Full_Update_Job_ID []int
	List_Systems       []string
}

type ScheduledJobs struct {
	AllJobs          []Job
	JobType          string
	Full_Update_Jobs Full_Update_Jobs
}

type Jobs_Patching struct {
	Patching             []interface{} `json:"Patching"`
	Reboot               []interface{} `json:"reboot_jobs"`
	JobcheckerEmails     []string      `json:"jobchecker_emails,omitempty"`
	JobcheckerTimeout    int           `json:"jobchecker_timeout,omitempty"`
	JobstartDelay        int           `json:"jobstart_delay,omitempty"`
	Offline_minions      []string      `json:"offline_minions,omitempty"`
	Disqualified_minions []string      `json:"btrfs_disqualified,omitempty"`
	T7user               string        `json:"t7user,omitempty"`
	Post_patching_file   string        `json:"post_patching_file,omitempty"`
	No_patch_execptions  []string      `json:"no_patch_execptions,omitempty"`
	Patch_level          string        `json:"patch_level,omitempty"`
	Post_patching        string        `json:"post_patching,omitempty"`
	Prep_patching        string        `json:"prep_patching,omitempty"`
	Full_Update_Job_ID   []interface{} `json:"Full_Update_Job_ID,omitempty"`
}

type Jobstatus struct {
	Pending              []Job    `json:"pending,omitempty"`
	Completed            []Job    `json:"completed,omitempty"`
	Failed               []Job    `json:"failee,omitempty"`
	Cancelled            []Job    `json:"cancelled,omitempty"`
	JobcheckerEmails     []string `json:"jobcheckeremails,omitempty"`
	JobcheckerTimeout    int      `json:"jobcheckerr timeout,omitempty"`
	JobstartDelay        int      `json:"job start delay,omitempty"`
	Offline_minions      []string `json:"offline minions,omitempty"`
	Disqualified_minions []string `json:"btrfs disqualified,omitempty"`
	No_patch_execptions  []string `json:"no_patch exceptions,omitempty"`
	T7user               string   `json:"t7user,omitempty"`
	JobStartTime         string   `json:"job start time,omitempty"`
	YamlFileName         string   `json:"completed yaml file,omitempty"`
	YamlFileName_Pending string   `json:"pending yaml file,omitempty"`
	YamlFileName_Failed  string   `json:"failed yaml file,omitempty"`
	YamlFileName_Full    string   `json:"full list yaml file,omitempty"`
	Reboot_List          string   `json:"reboot list,omitempty"`
	Reboot_SLS           string   `json:"reboot sls,omitempty"`
	Reboot_Command       string   `json:"reboot command,omitempty"`
	Post_patching_file   string   `json:"post patching file,omitempty"`
	Patch_level          string   `json:"patch level,omitempty"`
	Post_patching        string   `json:"post patch sls,omitempty"`
	Prep_patching        string   `json:"prep patch sls,omitempty"`
	JobType              string   `json:"job type,omitempty"`
}

type ListSystemInJobs struct {
	ListInProgressSystems ListSystemInJobs_Response
	ListFailedSystems     ListSystemInJobs_Response
	ListCompletedSystems  ListSystemInJobs_Response
}

/* type ListSystemInJobs_Response struct {
	Result []struct {
		Server_name  string    `xmlrpc:"server_name,omitempty"`
		Base_channel string    `xmlrpc:"base_channel,omitempty"`
		Server_id    int       `xmlrpc:"server_id,omitempty"`
		Timestamp    time.Time `xmlrpc:"timestamp,omitempty"`
	}
} */

type ListSystemInJobs_Response struct {
	Result []struct {
		Server_name  string
		Base_channel string
		Server_id    int
		Timestamp    time.Time
	}
}

type ListSystemInJobs_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	ActionId   int    `xmlrpc:"actionId"`
}
