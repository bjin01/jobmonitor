package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"net/http"
	"os"
	"time"

	"github.com/bjin01/jobmonitor/delete_systems"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/pkg_updates"
	"github.com/bjin01/jobmonitor/saltapi"
	"github.com/bjin01/jobmonitor/schedules"
	"github.com/bjin01/jobmonitor/spmigration"
	"github.com/gin-gonic/gin"
)

func init() {

}

type my_context struct {
	ctx      context.Context
	cancells context.CancelFunc
}

type keyType struct {
	name string
}

var time_counter_pkg_update_by_list int
var time_counter_pkg_update int

func main() {
	sumafile_path := flag.String("config", "/etc/salt/master.d/spacewalk.conf", "provide config file with SUMA login")
	api_interval := flag.Int("interval", 60, "SUMA API polling interval, default 60 seconds, no need to write s.")
	templates_dir := flag.String("templates", "/srv/jobmonitor", "provide directory name where the template files are stored.")
	port := flag.Int("port", 12345, "provide port number for the web server.")
	tls_cert := flag.String("tls_cert", "/etc/pki/tls/certs/spacewalk.crt", "provide tls certificate file path.")
	tls_key := flag.String("tls_key", "/etc/pki/tls/private/spacewalk.key", "provide tls key file path.")
	flag.Parse()

	var my_contexts []my_context
	SUMAConfig := GetConfig(*sumafile_path)
	logger.Infof("interval is: %v\n", *api_interval)

	if _, err := os.Stat(*templates_dir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		logger.Fatalf("templates directory missing: %s\n", *templates_dir)
	}
	templates := &email.Templates_Dir{Dir: *templates_dir}
	logger.Infof("templates_dir is: %s\n", templates.Dir)
	logger.Infof("port is: %v\n", *port)

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
				logger.Infof("Health check email recipients: %v\n", *emails_to)
			}

			if healthcheck_interval_field.IsValid() && b.Healthcheck_interval > 60 {
				*healthcheck_interval = b.Healthcheck_interval
			} else {
				*healthcheck_interval = 60
			}
			logger.Infof("Health check interval: %ds\n", *healthcheck_interval)
		}
		//error_counter := 0
		interval := time.Duration(*healthcheck_interval) * time.Second
		ticker := time.NewTicker(interval)
		// Call the API immediately
		err := performHealthCheck(SUMAConfig)
		if err != nil {
			*health = false
			logger.Infof("SUSE Manager initial health check failed. %v %s\n", *health, err)
		} else {
			*health = true
			logger.Infof("SUSE Manager health check passed. %v\n", *health)
		}

		// Start the loop to perform the API call periodically
		for range ticker.C {
			err = performHealthCheck(SUMAConfig)
			if err != nil {
				*health = false
				logger.Infof("SUSE Manager health check failed. %s\n", err)
			} else {
				*health = true
				logger.Infof("SUSE Manager health check passed. %v\n", *health)
			}

			/* subject := "SUSE Manager health check failed"
			message := fmt.Sprintf("SUSE Manager health check failed 5 times in serie.\n")
			if len(*emails_to) > 0 {
				email.Send_system_emails(*emails_to, subject, message)
			} else {
				logger.Infof("Alarm: SUSE Manager health check failed 5 times in row.")
			} */

		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.LoadHTMLGlob(fmt.Sprintf("%s/*", *templates_dir))

	static_dir_split := strings.Split(*templates_dir, "/")
	// trim out the last / part
	static_dir_list := static_dir_split[:len(static_dir_split)-1]
	// join the string back
	static_dir := strings.Join(static_dir_list, "/")
	static_dir = fmt.Sprintf("%s/static", static_dir)
	logger.Infoln("static_dir: ", static_dir)
	r.Static("/static", static_dir)

	r.GET("/web", func(c *gin.Context) {
		Myname := "Bo Jin"
		c.HTML(http.StatusOK, "web.html", Myname)
	})

	r.POST("/viewdb", func(c *gin.Context) {
		dbfile := c.PostForm("dbfile")
		fmt.Println("POST Request for dbfile: ", dbfile)
		// Store dbfile path in a cookie
		c.SetCookie("dbfile", dbfile, 3600, "/", "", false, true)

		db, err := sql.Open("sqlite3", dbfile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open database: %v", err)
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to query database: %v", err)
			return
		}
		defer rows.Close()

		//var tables []string
		/* html := `<!DOCTYPE html>
		<html>
		<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<body>`
		html += fmt.Sprintf("<p class=\"fs-1\">Tables from %s</p>", dbfile)
		html += `<div class="list-group">`
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				c.String(http.StatusInternalServerError, "Failed to scan row: %v", err)
				return
			}
			html += fmt.Sprintf("<a href=\"/table?name=%s\" class=\"list-group-item list-group-item-action\">%s</a>", table, table)
		}

		html += `</div>
		<a class="btn btn-primary" href="/web" role="button">Enter DB file</a>
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" crossorigin="anonymous"></script>
		</body></html>` */
		return_values := []string{}
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				c.String(http.StatusInternalServerError, "Failed to scan row: %v", err)
				return
			}
			return_values = append(return_values, table)
		}
		logger.Debugln("return_values: ", return_values)
		c.JSONP(http.StatusOK, return_values)
	})

	r.GET("/viewdb", func(c *gin.Context) {
		// Retrieve dbfile path from cookie
		dbfile, err := c.Cookie("dbfile")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get database file path: %v", err)
			return
		}

		db, err := sql.Open("sqlite3", dbfile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open database: %v", err)
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to query database: %v", err)
			return
		}
		defer rows.Close()

		html := `<!DOCTYPE html>
		<html>
		<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<body>`
		html += fmt.Sprintf("<p class=\"fs-1\">Tables from %s</p>", dbfile)
		html += `<div class="list-group">`

		//var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				c.String(http.StatusInternalServerError, "Failed to scan row: %v", err)
				return
			}
			html += fmt.Sprintf("<a href=\"/table?name=%s\" class=\"list-group-item list-group-item-action\">%s</a>", table, table)
		}

		html += `</div>
		<a class="btn btn-primary" href="/web" role="button">Enter DB file</a>
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" crossorigin="anonymous"></script>
		</body></html>`

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		//c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("Tables: %s", strings.Join(tables, ", "))))
	})

	r.GET("/table", func(c *gin.Context) {
		tableName := c.Query("name")

		dbfile, err := c.Cookie("dbfile")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get database file path: %v", err)
			return
		}
		db, err := sql.Open("sqlite3", dbfile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open database: %v", err)
			return
		}
		defer db.Close()

		logger.Debugln("Query table: ", tableName)

		query := fmt.Sprintf("SELECT * FROM %s", tableName)
		if tableName == "minion_data" {
			query = fmt.Sprintf("SELECT minion_name, minion_status, minion_remarks, migration_stage, migration_stage_status FROM %s ORDER BY Minion_Name", tableName)
		}

		rows, err := db.Query(query)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to query table: %v", err)
			return
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get columns: %v", err)
			return
		}

		var result []map[string]interface{}
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := 0; i < len(columns); i++ {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			rowResult := make(map[string]interface{})
			for i, colName := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				rowResult[colName] = v
				//fmt.Println(colName, v)
			}

			result = append(result, rowResult)
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to marshal result: %v", err)
			return
		}

		c.String(http.StatusOK, string(jsonData))
		/* columns, err := rows.Columns()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get columns: %v", err)
			return
		} */

		/* var tableContent []map[string]interface{}
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for rows.Next() {
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			row := make(map[string]interface{})
			for i, colName := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				row[colName] = v
			}
			tableContent = append(tableContent, row)
		}

		// Start the HTML table
		tableHTML := "<table class='table table-hover'><thead><tr>"

		// Add table headers
		for _, colName := range columns {
			tableHTML += fmt.Sprintf("<th scope=col>%s</th>", colName)
		}

		// End table headers row
		tableHTML += "</tr></thead><tbody>"

		// Add table rows
		for _, row := range tableContent {
			tableHTML += "<tr>"
			for _, colName := range columns {
				tableHTML += fmt.Sprintf("<td>%v</td>", row[colName])
			}
			tableHTML += "</tr>"
		}

		// End the HTML table
		tableHTML += "</tbody></table>"

		html := `<!DOCTYPE html>
		<html>
		<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<body>`
		// Add the HTML table to the HTML document
		html += fmt.Sprintf("<p class=\"fs-1\">Table %s from %s</p>", tableName, dbfile)
		html += tableHTML

		// Add the "Go Back" link to the HTML document
		html += `<a class="btn btn-primary" href="/viewdb" role="button">Back to table list</a>`

		// End the HTML document
		html += `
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" crossorigin="anonymous"></script>
		</body></html>`

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html)) */
	})

	r.POST("/delete_system", func(c *gin.Context) {
		/* minionName := c.PostForm("minion_name")
		authToken := c.GetHeader("Authentication-Token") */
		Setup_Logger("")
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

		c.JSON(http.StatusOK, data)
	})

	r.POST("/salt", func(c *gin.Context) {
		var saltdata saltapi.Salt_Data
		if err := c.ShouldBindJSON(&saltdata); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		logger.Fatalf("saltdata: %v\n", saltdata)
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
		//logger.Fatalf("SaltJob_Data: %v\n", saltdata.Jid)
		saltdata.Login()
		saltdata.Query_Jid()
		if saltdata.Token != "" {
			c.Data(http.StatusOK, "application/json; charset=utf-8", saltdata.Return)

		} else {
			c.JSON(http.StatusOK, gin.H{"error": "Authentication failed"})
		}

	})

	r.POST("/pkg_update", func(c *gin.Context) {
		var pkg_update_request_obj pkg_updates.Update_Groups

		if !*health {
			c.String(200, "Pkg Update will not start due to SUSE Manager health check failed. Please check the logs.")
			logger.Infof("Pkg Update will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}

		request_received := new(time.Time)
		*request_received = time.Now()
		if time_counter_pkg_update == 0 {
			if err := c.ShouldBindJSON(&pkg_update_request_obj); err != nil {

				c.AbortWithError(http.StatusBadRequest, err)
			}

			if !isValidAuthToken(pkg_update_request_obj.Token) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			Setup_Logger(pkg_update_request_obj.Logfile)
			ctx, cancel := context.WithCancel(context.Background())

			ctx_timestemp := fmt.Sprintf("%s_%d", pkg_update_request_obj.T7User, time.Now().Nanosecond())

			var key = keyType{"ctx_timestemp"}

			ctx = context.WithValue(ctx, key, ctx_timestemp)
			context_temp := my_context{ctx, cancel}
			my_contexts = append(my_contexts, context_temp)
			//context_value := ctx.Value("ctx_timestemp")
			context_value := ctx.Value(key)
			pkg_update_request_obj.Ctx_ID = ctx_timestemp
			logger.Infof("------------------context value: %v\n", context_value.(string))
			logger.Infoln("Length of context: ", len(my_contexts))
			go Pkg_update_groups_lookup(ctx, SUMAConfig, &pkg_update_request_obj, templates, health)
			c.String(http.StatusOK, fmt.Sprintf("Targeting %v for Package Updates & SP Migration through SUSE Manager.", pkg_update_request_obj.Groups))
			//logger.Infof("request data %v for Package Update through SUSE Manager.\n", pkg_update_request_obj)
		} else {
			c.String(http.StatusOK, "Package Updates & SP Migration POST request repeated too soon. Please wait for 20 seconds.")
			logger.Infof("Package Updates & SP Migrationrequest repeated too soon. Please wait for 20 seconds.")
		}

		go func() {

			protect_time_window := request_received.Add(time.Duration(20) * time.Second)
			for time.Now().Before(protect_time_window) {
				time_counter_pkg_update += 1
				time.Sleep(time.Second * 1)
			}
			time_counter_pkg_update = 0
		}()
	})

	r.POST("/pkg_update_by_list", func(c *gin.Context) {
		var pkg_update_request_obj pkg_updates.Update_Groups

		if !*health {
			c.String(200, "Pkg Update will not start due to SUSE Manager health check failed. Please check the logs.")
			logger.Infof("Pkg Update will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}

		request_received := new(time.Time)
		*request_received = time.Now()

		if err := c.ShouldBindJSON(&pkg_update_request_obj); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		if !isValidAuthToken(pkg_update_request_obj.Token) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//logger.Fatalf("pkg_update_request_obj %+v\n", pkg_update_request_obj)

		/* go Pkg_update_by_list(SUMAConfig, &pkg_update_request_obj, templates, health) */

		Setup_Logger(pkg_update_request_obj.Logfile)
		if time_counter_pkg_update_by_list == 0 {
			go Pkg_update_by_list(SUMAConfig, &pkg_update_request_obj, templates, health)
			c.String(http.StatusOK, "Targeting %v for Package Updates & SP Migration through SUSE Manager.", pkg_update_request_obj.Minions_to_add)
			logger.Infof("request data for Package Update through SUSE Manager received.\n")
		} else {
			c.String(http.StatusOK, "Pkg Update POST request repeated too soon. Please wait for 20 seconds.")
			logger.Infof("Pkg Update POST request repeated too soon. Please wait for 20 seconds.")
		}

		go func() {

			protect_time_window := request_received.Add(time.Duration(20) * time.Second)
			for time.Now().Before(protect_time_window) {
				time_counter_pkg_update_by_list += 1
				time.Sleep(time.Second * 1)
			}
			time_counter_pkg_update_by_list = 0
		}()

	})

	r.GET("/cancell_pkg_update", func(c *gin.Context) {
		ctx_id := c.Query("ctx_id")
		if ctx_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'ctx_id' parameter"})
			return
		}

		var key = keyType{"ctx_timestemp"}

		for i, ctx := range my_contexts {
			context_value := ctx.ctx.Value(key)
			logger.Infof("------------------context value: %v\n", context_value.(string))
			if context_value.(string) == ctx_id {
				logger.Infof("context value found: %v\n", context_value.(string))
				ctx.cancells()
				my_contexts = append(my_contexts[:i], my_contexts[i+1:]...)
				logger.Infof("Length of context: %v\n", len(my_contexts))
				c.String(http.StatusOK, fmt.Sprintf("Context %s is cancelled.", ctx_id))
				return
			}
		}
	})

	r.GET("/pkg_update", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'filename' parameter"})
			return
		}

		minion_name := c.Query("minion_name")
		if minion_name != "" {
			result, err := Pkg_update_get_minion_from_db(filename, minion_name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.IndentedJSON(http.StatusOK, result)
			return
		} else {
			result := Pkg_update_groups_lookup_from_file(filename)
			c.IndentedJSON(http.StatusOK, result)
		}

	})

	r.POST("/spmigration", func(c *gin.Context) {
		var spmigrationRequestObj spmigration.Migration_Groups

		if !*health {
			c.String(200, "SPMigration will not start due to SUSE Manager health check failed. Please check the logs.")
			logger.Infof("SPMigration will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}
		if err := c.ShouldBindJSON(&spmigrationRequestObj); err != nil {

			c.AbortWithError(http.StatusBadRequest, err)
		}

		if !isValidAuthToken(spmigrationRequestObj.Token) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//logger.Fatalf("spmigrationRequestObj %+v\n", spmigrationRequestObj)

		go groups_lookup(SUMAConfig, &spmigrationRequestObj, templates, health)
		c.String(http.StatusOK, fmt.Sprintf("Targeting %v for SP Migration through SUSE Manager.", spmigrationRequestObj.Groups))
		//logger.Infof("request data %v for SP Migration through SUSE Manager.\n", spmigrationRequestObj)

	})

	r.GET("/query_jobchecker", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'filename' parameter"})
			return
		}

		data, err := readJSONFile_patching(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/jobchecker", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		//logger.Fatalf("cCp %+v\n", cCp.Request.Body)
		var instance_jobs_patching schedules.Jobs_Patching
		var alljobs schedules.ScheduledJobs
		var full_update_jobs schedules.Full_Update_Jobs

		if !*health {
			c.String(200, "Jobchecker will not start due to SUSE Manager health check failed. Please check the logs.")
			logger.Infof("Jobchecker will not start due to SUSE Manager health check failed. Please check the logs.")
			return
		}

		if err := cCp.ShouldBindJSON(&instance_jobs_patching); err != nil {
			cCp.AbortWithError(http.StatusBadRequest, err)
		}

		for _, elem := range instance_jobs_patching.Patching {
			if instance_jobs_patching.Full_Update_Job_ID != nil {
				for _, elem := range instance_jobs_patching.Full_Update_Job_ID {
					for k, v := range elem.(map[string]interface{}) {
						//logger.Fatalf("k: %v, elem: %v\n", k, elem)
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
				logger.Infof("Bundle full_update_jobs %+v\n", full_update_jobs)
				alljobs.Full_Update_Jobs = full_update_jobs
			}
			for x, y := range elem.(map[string]interface{}) {
				var instance_jobs_patching schedules.Job
				//logger.Infof("%s:\n", x)
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
				logger.Infof("instance_jobs_patching %+v\n", instance_jobs_patching)
			}
		}

		go Jobmonitor(SUMAConfig, alljobs, instance_jobs_patching, templates, health)
		c.String(200, "Jobchecker task started.")
	})
	logger.Infof("/jobckecker API is listening and serving HTTP on :%d", *port)
	// Listen and serve on given port
	//r.Run(fmt.Sprintf(":%d", *port))

	// Run the server with TLS.
	// Replace 'cert.pem' and 'key.pem' with the paths to your actual certificate and key files.
	err := r.RunTLS(fmt.Sprintf(":%d", *port), *tls_cert, *tls_key)
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}

}
