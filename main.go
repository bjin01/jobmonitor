package main

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/bjin01/jobmonitor/delete_systems"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/saltapi"
	"github.com/bjin01/jobmonitor/schedules"
	"github.com/bjin01/jobmonitor/spmigration"
	"github.com/gin-gonic/gin"
)

func init() {

}

func main() {
	sumafile_path := flag.String("config", "/etc/salt/master.d/spacewalk.conf", "provide config file with SUMA login")
	api_interval := flag.Int("interval", 60, "SUMA API polling interval, default 60 seconds, no need to write s.")
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

	health := new(bool)
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
		//error_counter := 0
		interval := time.Duration(*healthcheck_interval) * time.Second
		ticker := time.NewTicker(interval)
		// Call the API immediately
		err := performHealthCheck(SUMAConfig)
		if err != nil {
			*health = false
			log.Printf("SUSE Manager initial health check failed. %v %s\n", *health, err)
		} else {
			*health = true
			log.Printf("SUSE Manager health check passed. %v\n", *health)
		}

		// Start the loop to perform the API call periodically
		for range ticker.C {
			err = performHealthCheck(SUMAConfig)
			if err != nil {
				*health = false
				log.Printf("SUSE Manager health check failed. %s\n", err)
			} else {
				*health = true
				log.Printf("SUSE Manager health check passed. %v\n", *health)
			}

			/* subject := "SUSE Manager health check failed"
			message := fmt.Sprintf("SUSE Manager health check failed 5 times in serie.\n")
			if len(*emails_to) > 0 {
				email.Send_system_emails(*emails_to, subject, message)
			} else {
				log.Println("Alarm: SUSE Manager health check failed 5 times in row.")
			} */

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

		go Delete_System(SUMAConfig, &deleteSystemRequestObj)
		c.String(http.StatusOK, fmt.Sprintf("System (%s) delete request is sent to SUSE Manager.", deleteSystemRequestObj.MinionName))
	})

	r.GET("/query_spmigration", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'filename' parameter"})
			return
		}

		data, err := readJSONFile(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data.MinionList)
	})

	r.POST("/salt", func(c *gin.Context) {
		var saltdata saltapi.Salt_Data
		if err := c.ShouldBindJSON(&saltdata); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		fmt.Printf("saltdata: %v\n", saltdata)
		saltdata.Login()
		saltdata.Run()
		if saltdata.Token != "" {
			c.Data(http.StatusOK, "application/json; charset=utf-8", saltdata.Return)

		} else {
			c.JSON(http.StatusOK, gin.H{"error": "Authentication failed"})
		}

	})

	r.POST("/saltjob", func(c *gin.Context) {
		var saltdata saltapi.Salt_Data
		if err := c.ShouldBindJSON(&saltdata); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		fmt.Printf("SaltJob_Data: %v\n", saltdata.Jid)
		saltdata.Login()
		saltdata.Query_Jid()
		if saltdata.Token != "" {
			c.Data(http.StatusOK, "application/json; charset=utf-8", saltdata.Return)

		} else {
			c.JSON(http.StatusOK, gin.H{"error": "Authentication failed"})
		}

	})

	r.POST("/spmigration", func(c *gin.Context) {
		var spmigrationRequestObj spmigration.Migration_Groups

		if *health == false {
			c.String(200, "SPMigration will not start due to SUSE Manager health check failed. Please check the logs.")
			log.Printf("SPMigration will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}

		if err := c.ShouldBindJSON(&spmigrationRequestObj); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		if !isValidAuthToken(spmigrationRequestObj.Token) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//fmt.Printf("spmigrationRequestObj %+v\n", spmigrationRequestObj)

		go groups_lookup(SUMAConfig, &spmigrationRequestObj, templates, health)
		c.String(http.StatusOK, fmt.Sprintf("Targeting %v for SP Migration through SUSE Manager.", spmigrationRequestObj.Groups))
		//log.Printf("request data %v for SP Migration through SUSE Manager.\n", spmigrationRequestObj)

	})

	r.POST("/jobchecker", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		//fmt.Printf("cCp %+v\n", cCp.Request.Body)
		var instance_jobs_patching schedules.Jobs_Patching
		var alljobs schedules.ScheduledJobs
		var full_update_jobs schedules.Full_Update_Jobs

		if *health == false {
			c.String(200, "Jobchecker will not start due to SUSE Manager health check failed. Please check the logs.")
			log.Printf("Jobchecker will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}

		if err := cCp.ShouldBindJSON(&instance_jobs_patching); err != nil {
			cCp.AbortWithError(http.StatusBadRequest, err)
		}

		for _, elem := range instance_jobs_patching.Patching {
			if instance_jobs_patching.Full_Update_Job_ID != nil {
				for _, elem := range instance_jobs_patching.Full_Update_Job_ID {
					for k, v := range elem.(map[string]interface{}) {
						//fmt.Printf("k: %v, elem: %v\n", k, elem)
						jobid, _ := strconv.Atoi(k)
						if len(full_update_jobs.Full_Update_Job_ID) == 0 {
							full_update_jobs.Full_Update_Job_ID = append(full_update_jobs.Full_Update_Job_ID, jobid)
						} else {
							for j := range full_update_jobs.Full_Update_Job_ID {
								if full_update_jobs.Full_Update_Job_ID[j] != jobid {
									full_update_jobs.Full_Update_Job_ID = append(full_update_jobs.Full_Update_Job_ID,
										jobid)
								}
							}
						}

						for _, v := range v.([]interface{}) {
							full_update_jobs.List_Systems = append(full_update_jobs.List_Systems, v.(string))
						}
					}
				}
				log.Printf("Bundle full_update_jobs %+v\n", full_update_jobs)
				alljobs.Full_Update_Jobs = full_update_jobs
			}
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

		go Jobmonitor(SUMAConfig, alljobs, instance_jobs_patching, templates, health)
		c.String(200, "Jobchecker task started.")
	})
	log.Default().Println("/jobckecker API is listening and serving HTTP on :12345")
	// Listen and serve on 0.0.0.0:12345
	r.Run(":12345")

}
