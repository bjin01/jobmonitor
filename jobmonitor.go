package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bjin01/jobmonitor/auth"
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
}

type SUMAConfig struct {
	SUMA map[string]struct {
		User     string `yaml:"username"`
		Password string `yaml:"password"`
		Logfile  string `yaml:"logfile"`
	} `yaml:"suma_api"`
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

func Jobmonitor(SUMAConfig *SUMAConfig, alljobs schedules.ScheduledJobs, instance_jobs_patching schedules.Jobs_Patching) {
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

	for _, job := range alljobs.AllJobs {
		log.Printf("Host: %s \tJob-ID: %d\n", job.Hostname, job.JobID)
	}

	if len(alljobs.AllJobs) != 0 {

		deadline := time.Now().Add(time.Duration(jobstatus_result.JobcheckerTimeout) * time.Minute)
		Jobstart_starttime := time.Now().Add(time.Duration(jobstatus_result.JobstartDelay) * time.Minute)
		jobstatus_result.JobStartTime = Jobstart_starttime.Format(time.RFC822Z)
		jobstatus_result.YamlFileName = fmt.Sprintf("completd_%s_%s", jobstatus_result.T7user, Jobstart_starttime.Format("20060102150405"))
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

					jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
					if err != nil {
						log.Default().Printf("ERROR: reboot list: %s\n", err)
					}
					jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result)
					if err != nil {
						log.Default().Printf("ERROR: reboot list: %s\n", err)
					}
					email.Sendit(jobstatus_result)
					break begin
				}
			}
			time.Sleep(10 * time.Second)
		}
		if len(jobstatus_result.Pending) > 0 {
			jobstatus_result.Reboot_List, err = email.WriteYaml(jobstatus_result)
			if err != nil {
				log.Default().Printf("reboot list: %s\n", err)
			}
			jobstatus_result.Reboot_SLS, err = email.Write_SLS(jobstatus_result)
			if err != nil {
				log.Default().Printf("ERROR: reboot list: %s\n", err)
			}
			email.Sendit(jobstatus_result)
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

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatalln("No config file for SUMA provided.")
	}
	if len(os.Args) <= 2 {
		log.Fatalf("Not enough arguments given. Two args needed. %+v\n", os.Args)
	}
	config_file := os.Args[1:]
	SUMAConfig := GetConfig(config_file[0])

	interval, err := strconv.Atoi(config_file[1])
	if err != nil {
		log.Printf("Failed to parse interval into integer: %v\n", err.Error())
	}
	log.Printf("interval is: %d\n", interval)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

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

		go Jobmonitor(SUMAConfig, alljobs, instance_jobs_patching)
		c.String(200, "Jobchecker task started.")
	})
	log.Default().Println("/jobckecker API is listening and serving HTTP on :12345")
	// Listen and serve on 0.0.0.0:12345
	r.Run(":12345")

}
