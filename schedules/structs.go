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
	Hostname   string
	ServerID   int
	JobID      int
	Masterplan string
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
	Pending              []Job
	Completed            []Job
	Failed               []Job
	Cancelled            []Job
	JobcheckerEmails     []string
	JobcheckerTimeout    int
	JobstartDelay        int
	Offline_minions      []string
	Disqualified_minions []string
	No_patch_execptions  []string
	T7user               string
	JobStartTime         string
	YamlFileName         string
	YamlFileName_Pending string
	YamlFileName_Failed  string
	Reboot_List          string
	Reboot_SLS           string
	Reboot_Command       string
	Post_patching_file   string
	Patch_level          string
	Post_patching        string
	Prep_patching        string
	JobType              string
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
