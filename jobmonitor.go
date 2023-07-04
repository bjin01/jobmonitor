package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"reflect"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/delete_systems"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/request"
	"github.com/bjin01/jobmonitor/schedules"
	"github.com/fernet/fernet-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Sumaconf struct {
	Server   string
	User     string
	Password string
	Email_to []string
}

type SUMAConfig struct {
	SUMA map[string]struct {
		User                 string   `yaml:"username"`
		Password             string   `yaml:"password"`
		Logfile              string   `yaml:"logfile"`
		Email_to             []string `yaml:"email_to"`
		Healthcheck_interval int      `yaml:"healthcheck_interval"`
		Healthcheck_email_to []string `yaml:"healthcheck_email"`
	} `yaml:"suma_api"`
}

func init() {

}

func GetConfig(file string) *SUMAConfig {
	// Read the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	// Create a struct to hold the YAML data
	var config SUMAConfig

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)

	}

	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
	}

	return &config
}

func Decrypt(key string, cryptoText string) string {
	k := fernet.MustDecodeKeys(key)
	/* tok, err := fernet.EncryptAndSign([]byte(cryptoText), k[0])
	if err != nil {
		panic(err)
	} */
	msg := fernet.VerifyAndDecrypt([]byte(cryptoText), 0, k)
	//fmt.Println(string(msg))

	return fmt.Sprintf("%s", msg)
}

func Delete_System(SUMAConfig *SUMAConfig, deleteSystemdata *delete_systems.DeleteSystemRequest) {
	fmt.Printf("deleteSystemdata %s\n", deleteSystemdata.MinionName)
	var sumaconf Sumaconf
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
	}
	for a, b := range SUMAConfig.SUMA {
		sumaconf.Server = a
		b.Password = Decrypt(key, b.Password)
		sumaconf.Password = b.Password
		sumaconf.User = b.User
		if len(b.Email_to) > 0 {
			sumaconf.Email_to = b.Email_to
		}
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sessionkey: %s\n", SessionKey.Sessionkey)
	fmt.Printf("Deleting System in SUMA: %s\n", deleteSystemdata.MinionName)
	err = delete_systems.Delete_System(SessionKey, deleteSystemdata, sumaconf.Email_to)
	if err != nil {
		log.Fatal(err)
	}

}

func Jobmonitor(SUMAConfig *SUMAConfig, alljobs schedules.ScheduledJobs,
	instance_jobs_patching schedules.Jobs_Patching, templates_dir *email.Templates_Dir) {
	//key := "R2bfp223Qsk-pX970Jw8tyJUChT4-e2J8anZ4G4n4IM="
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
	}

	var sumaconf Sumaconf
	for a, b := range SUMAConfig.SUMA {
		sumaconf.Server = a
		b.Password = Decrypt(key, b.Password)
		sumaconf.Password = b.Password
		sumaconf.User = b.User
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		log.Fatal(err)
	}

	jobstatus_result := new(schedules.Jobstatus)

	switch alljobs.JobType {
	case "patching":
		jobstatus_result.JobType = "patching"
	case "reboot":
		jobstatus_result.JobType = "reboot"
	default:
		jobstatus_result.JobType = "patching"
	}

	if len(instance_jobs_patching.JobcheckerEmails) != 0 {
		jobstatus_result.JobcheckerEmails = instance_jobs_patching.JobcheckerEmails
		log.Printf("JobcheckerEmails : %v\n", jobstatus_result.JobcheckerEmails)
	} else {
		jobstatus_result.JobcheckerEmails = []string{}
		log.Printf("JobcheckerEmails : %v\n", jobstatus_result.JobcheckerEmails)
	}

	if instance_jobs_patching.JobcheckerTimeout != 0 {
		jobstatus_result.JobcheckerTimeout = instance_jobs_patching.JobcheckerTimeout
		log.Printf("JobcheckerTimeout : %v\n", jobstatus_result.JobcheckerTimeout)
	} else {
		jobstatus_result.JobcheckerTimeout = 2
		log.Printf("JobcheckerTimeout : %v\n", jobstatus_result.JobcheckerTimeout)
	}

	if len(instance_jobs_patching.Offline_minions) != 0 {
		jobstatus_result.Offline_minions = instance_jobs_patching.Offline_minions
		log.Printf("Offline_minions : %v\n", jobstatus_result.Offline_minions)
	} else {
		jobstatus_result.Offline_minions = []string{}
		log.Printf("Offline_minions : %v\n", jobstatus_result.Offline_minions)
	}

	if len(instance_jobs_patching.Disqualified_minions) != 0 {
		jobstatus_result.Disqualified_minions = instance_jobs_patching.Disqualified_minions
		log.Printf("Disqualified_minions : %v\n", jobstatus_result.Disqualified_minions)
	} else {
		jobstatus_result.Disqualified_minions = []string{}
		log.Printf("Disqualified_minions : %v\n", jobstatus_result.Disqualified_minions)
	}

	if len(instance_jobs_patching.No_patch_execptions) != 0 {
		jobstatus_result.No_patch_execptions = instance_jobs_patching.No_patch_execptions
		log.Printf("No_patch_execptions : %v\n", jobstatus_result.No_patch_execptions)
	} else {
		jobstatus_result.No_patch_execptions = []string{}
		log.Printf("No_patch_execptions : %v\n", jobstatus_result.No_patch_execptions)
	}

	if instance_jobs_patching.JobstartDelay != 0 {
		jobstatus_result.JobstartDelay = instance_jobs_patching.JobstartDelay
		log.Printf("JobstartDelay : %v\n", jobstatus_result.JobstartDelay)
	} else {
		jobstatus_result.JobstartDelay = 1
		log.Printf("JobstartDelay : %v\n", jobstatus_result.JobstartDelay)
	}

	if instance_jobs_patching.T7user != "" {
		jobstatus_result.T7user = instance_jobs_patching.T7user
		log.Printf("T7user : %v\n", jobstatus_result.T7user)
	} else {
		jobstatus_result.T7user = "unknown"
		log.Printf("T7user : %v\n", jobstatus_result.T7user)
	}

	if instance_jobs_patching.Post_patching_file != "" {
		jobstatus_result.Post_patching_file = instance_jobs_patching.Post_patching_file
		log.Printf("Post_patching_file : %v\n", jobstatus_result.Post_patching_file)
	} else {
		jobstatus_result.Post_patching_file = ""
		log.Printf("Post_patching_file : %v\n", jobstatus_result.Post_patching_file)
	}

	if instance_jobs_patching.Post_patching != "" {
		jobstatus_result.Post_patching = instance_jobs_patching.Post_patching
		log.Printf("Post_patching state : %v\n", jobstatus_result.Post_patching)
	} else {
		jobstatus_result.Post_patching = ""
		log.Printf("Post_patching state : %v\n", jobstatus_result.Post_patching)
	}

	if instance_jobs_patching.Prep_patching != "" {
		jobstatus_result.Prep_patching = instance_jobs_patching.Prep_patching
		log.Printf("Prep_patching state : %v\n", jobstatus_result.Prep_patching)
	} else {
		jobstatus_result.Prep_patching = ""
		log.Printf("Prep_patching state : %v\n", jobstatus_result.Prep_patching)
	}

	if instance_jobs_patching.Patch_level != "" {
		jobstatus_result.Patch_level = instance_jobs_patching.Patch_level
		log.Printf("Patch_level : %v\n", jobstatus_result.Patch_level)
	} else {
		jobstatus_result.Patch_level = ""
		log.Printf("Patch_level : %v\n", jobstatus_result.Patch_level)
	}

	for _, job := range alljobs.AllJobs {
		log.Printf("Host: %s \tJob-ID: %d\n", job.Hostname, job.JobID)
	}

	if len(alljobs.AllJobs) != 0 {

		deadline := time.Now().Add(time.Duration(jobstatus_result.JobcheckerTimeout) * time.Minute)
		Jobstart_starttime := time.Now().Add(time.Duration(jobstatus_result.JobstartDelay) * time.Minute)
		jobstatus_result.JobStartTime = Jobstart_starttime.Format(time.RFC822Z)
		jobstatus_result.YamlFileName = fmt.Sprintf("completed_%s_%s", jobstatus_result.T7user, Jobstart_starttime.Format("20060102150405"))
	begin:
		for time.Now().Before(deadline) {

			log.Printf("Looping every minute. Deadline is %+v\n", deadline)
			log.Printf("Jobcheck will start at %+v\n", Jobstart_starttime)
			for time.Now().After(Jobstart_starttime) {
				jobstatus_result.Pending = []schedules.Job{}
				jobstatus_result.Failed = []schedules.Job{}
				jobstatus_result.Completed = []schedules.Job{}
				jobstatus_result.Cancelled = []schedules.Job{}
				jobstatus_result.Compare(SessionKey, alljobs.AllJobs)
				log.Printf("Jobstatus result: %+v\n", jobstatus_result)
				time.Sleep(10 * time.Second)
				if len(jobstatus_result.Pending) == 0 {
					log.Printf("No more pending Jobs. Exit loop. Email sent.")
					if jobstatus_result.JobType == "patching" {
						jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
						if err != nil {
							log.Default().Printf("ERROR: reboot list: %s\n", err)
						}
						jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result, templates_dir)
						if err != nil {
							log.Default().Printf("ERROR: reboot sls: %s\n", err)
						}
					}

					email.Sendit(jobstatus_result, templates_dir)
					break begin
				}
			}
			time.Sleep(10 * time.Second)
		}
		if len(jobstatus_result.Pending) > 0 {
			if jobstatus_result.JobType == "patching" {
				jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
				if err != nil {
					log.Default().Printf("ERROR: reboot list: %s\n", err)
				}
				jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result, templates_dir)
				if err != nil {
					log.Default().Printf("ERROR: reboot sls: %s\n", err)
				}
			}
			email.Sendit(jobstatus_result, templates_dir)
		}
	} else {
		log.Println("No Patch Jobs found.")
	}

	log.Println("Jobchecker timeout reached or not more jobs in pending status. Email sent.")
	err = auth.Logout("auth.logout", *SessionKey)
	if err != nil {
		log.Fatal(err)
	}
}

func isValidAuthToken(token string) bool {
	if token == os.Getenv("SUMAKEY") {
		return true
	} else {
		return false
	}

}

func performHealthCheck(sumaconfig *SUMAConfig) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	sumaserver := new(string)
	for a := range sumaconfig.SUMA {
		*sumaserver = a
	}

	sumaurl := fmt.Sprintf("https://%s/rhn/manager/api/api/systemVersion", *sumaserver)
	//log.Printf("suma url health check: %s\n", sumaurl)

	resp, err := client.Get(sumaurl)
	if err != nil {
		log.Println("SUMA Health check - API call failed:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code not OK, %d\n", resp.StatusCode)
	} else {
		log.Printf("Health check status: OK, %d\n", resp.StatusCode)
	}
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("Response: %s\n", string(responseBody)) */
	return nil
}

func main() {
	sumafile_path := flag.String("config", "/etc/salt/master.d/spacewalk.conf", "provide config file with SUMA login")
	api_interval := flag.Int("interval", 10, "SUMA API polling interval, default 10seconds, no need to write s.")
	templates_dir := flag.String("templates", "/srv/jobmonitor", "provide directory name where the template files are stored.")
	flag.Parse()

	SUMAConfig := GetConfig(*sumafile_path)
	log.Printf("interval is: %v\n", *api_interval)

	if _, err := os.Stat(*templates_dir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Fatalf("templates directory missing: %s\n", *templates_dir)
	}
	templates := &email.Templates_Dir{Dir: *templates_dir}
	log.Printf("templates_dir is: %s\n", templates.Dir)

	go func() {
		// Create a ticker with the desired interval
		healthcheck_interval := new(int)
		emails_to := new([]string)
		for _, b := range SUMAConfig.SUMA {
			value := reflect.ValueOf(b)
			healthcheck_email_to_field := value.FieldByName("Healthcheck_email_to")
			healthcheck_interval_field := value.FieldByName("Healthcheck_interval")
			if healthcheck_email_to_field.IsValid() && len(b.Healthcheck_email_to) > 0 {
				*emails_to = b.Healthcheck_email_to
				log.Printf("Health check email recipients: %v\n", *emails_to)
			}

			if healthcheck_interval_field.IsValid() && b.Healthcheck_interval > 60 {
				*healthcheck_interval = b.Healthcheck_interval
			} else {
				*healthcheck_interval = 60
			}
			log.Printf("Health check interval: %ds\n", *healthcheck_interval)
		}
		error_counter := 0
		interval := time.Duration(*healthcheck_interval) * time.Second
		ticker := time.NewTicker(interval)
		// Call the API immediately
		err := performHealthCheck(SUMAConfig)
		if err != nil {
			log.Fatalf("SUSE Manager initial health check failed. %s\n", err)
		}

		// Start the loop to perform the API call periodically
		for range ticker.C {
			err = performHealthCheck(SUMAConfig)
			if err != nil {
				error_counter += 1
			}
			if error_counter == 5 {
				subject := "SUSE Manager health check failed"
				message := fmt.Sprintf("SUSE Manager health check failed 5 times in serie.\n")
				if len(*emails_to) > 0 {
					email.Send_system_emails(*emails_to, subject, message)
				} else {
					log.Println("Alarm: SUSE Manager health check failed 5 times in row.")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/delete_system", func(c *gin.Context) {
		/* minionName := c.PostForm("minion_name")
		authToken := c.GetHeader("Authentication-Token") */
		var deleteSystemRequestObj delete_systems.DeleteSystemRequest
		if err := c.ShouldBindJSON(&deleteSystemRequestObj); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		// Perform authentication/token validation
		if !isValidAuthToken(deleteSystemRequestObj.Token) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Perform deletion of the system with the provided minion_name
		/* if err := deleteSystem(minionName); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} */

		go Delete_System(SUMAConfig, &deleteSystemRequestObj)
		c.String(http.StatusOK, fmt.Sprintf("System (%s) delete request is sent to SUSE Manager.", deleteSystemRequestObj.MinionName))
	})

	r.POST("/jobchecker", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		var instance_jobs_patching schedules.Jobs_Patching
		var alljobs schedules.ScheduledJobs

		if err := cCp.ShouldBindJSON(&instance_jobs_patching); err != nil {
			cCp.AbortWithError(http.StatusBadRequest, err)
		}

		for _, elem := range instance_jobs_patching.Patching {
			for x, y := range elem.(map[string]interface{}) {
				var instance_jobs_patching schedules.Job
				//log.Printf("%s:\n", x)
				instance_jobs_patching.Hostname = x
				for k, v := range y.(map[string]interface{}) {
					switch k {
					case "Patch Job ID is":
						instance_jobs_patching.JobID = int(v.(float64))
						alljobs.JobType = "patching"
					case "Reboot Job ID is":
						instance_jobs_patching.JobID = int(v.(float64))
						alljobs.JobType = "reboot"
					case "masterplan":
						instance_jobs_patching.Masterplan = v.(string)
					default:
						continue
					}
				}

				alljobs.AllJobs = append(alljobs.AllJobs, instance_jobs_patching)
				log.Printf("instance_jobs_patching %+v\n", instance_jobs_patching)
			}
		}

		go Jobmonitor(SUMAConfig, alljobs, instance_jobs_patching, templates)
		c.String(200, "Jobchecker task started.")
	})
	log.Default().Println("/jobckecker API is listening and serving HTTP on :12345")
	// Listen and serve on 0.0.0.0:12345
	r.Run(":12345")

}
