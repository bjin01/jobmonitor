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
	Completed ListCompletedJobs
	Failed    ListFailedJobs
	Pending   ListPendingJobs
}

type Job struct {
	Hostname   string
	JobID      int
	Masterplan string
}

type ScheduledJobs struct {
	AllJobs []Job
}

type Jobs_Patching struct {
	Patching             []interface{} `json:"Patching"`
	JobcheckerEmails     []string      `json:"jobchecker_emails,omitempty"`
	JobcheckerTimeout    int           `json:"jobchecker_timeout,omitempty"`
	JobstartDelay        int           `json:"jobstart_delay,omitempty"`
	Offline_minions      []string      `json:"offline_minions,omitempty"`
	Disqualified_minions []string      `json:"btrfs_disqualified,omitempty"`
	T7user               string        `json:"t7user,omitempty"`
	Post_patching_file   string        `json:"post_patching_file,omitempty"`
	No_patch_execptions  []string      `json:"no_patch_execptions,omitempty"`
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
	Reboot_List          string
	Reboot_SLS           string
	Reboot_Command       string
	Post_patching_file   string
}
