package main

import (
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/pkg_updates"
	"github.com/bjin01/jobmonitor/request"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Pkg_update_by_list(SUMAConfig *SUMAConfig, groupsdata *pkg_updates.Update_Groups,
	email_template_dir *email.Templates_Dir, health *bool) {

	if health != nil {
		if !*health {
			logger.Infof("Health check failed. Skipping groups lookup.")
			return
		}
	}

	//logger.Info("SP Migration input data %v\n", groupsdata)
	var sumaconf Sumaconf
	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		logger.Infof("SUMAKEY is not set. This might cause error for password decryption.")
	}
	for a, b := range SUMAConfig.SUMA {
		sumaconf.Server = a
		b.Password = Decrypt(key, b.Password)
		sumaconf.Password = b.Password
		sumaconf.User = b.User
		if len(b.Email_to) > 0 {
			sumaconf.Email_to = b.Email_to
		} else {
			sumaconf.Email_to = nil // or a suitable default value
		}
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		logger.Fatalln(err)
	}
	//email_template_directory_string := fmt.Sprintf("%s", email_template_dir.Dir)

	logger.Infof("Use sqlite database: %s\n", groupsdata.Sqlite_db)
	db, err := gorm.Open(gorm.Dialector(&sqlite.Dialector{DSN: groupsdata.Sqlite_db}),
		&gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	// Create the DB schema
	/* db.AutoMigrate(&pkg_updates.Workflow_Step{})
	db.AutoMigrate(&pkg_updates.Jobchecker_Email{})
	db.AutoMigrate(&pkg_updates.Group{})
	db.AutoMigrate(&pkg_updates.OptionalChannels{})
	db.AutoMigrate(&pkg_updates.Minion_Data{})

	var workflow_steps []pkg_updates.Workflow_Step
	for _, g := range groupsdata.Workflow {
		new_workflow := new(pkg_updates.Workflow_Step)
		for name, step := range g {
			new_workflow.Name = name
			new_workflow.Order = step

			workflow_steps = append(workflow_steps, *new_workflow)
			//fmt.Printf("-----------new_workflow: %+v\n", new_workflow)
			result := db.FirstOrCreate(&new_workflow, new_workflow)
			if result.RowsAffected > 0 {
				logger.Infof("Created workflow step %s - %d\n", name, result.RowsAffected)
			} else {
				db.Model(&new_workflow).Where("Name = ?", name).Update("Order", step)
				logger.Infof("Workflow Step %s already exists\n", name)
			}
		}
	} */

	/* wf, err := pkg_updates.Get_Workflow_From_DB(db)
	if err != nil {
		logger.Errorf("failed to get workflow from DB.")
		return
	}

	for _, w := range wf {
		logger.Infof("Workflow Step: %s - %d\n", w.Name, w.Order)
	} */

	var minion_list []pkg_updates.Minion_Data
	for _, g := range groupsdata.Minions_to_add {

		minion_id := pkg_updates.Get_SID(SessionKey, g)

		//it is important to delete the notes from system in suma here
		pkg_updates.Delete_Notes(SessionKey, minion_id)

		if minion_id != 0 {
			minion_list = append(minion_list, pkg_updates.Minion_Data{Minion_Name: g, Minion_ID: minion_id, Minion_Status: "Offline"})
			logger.Infof("SystemID for %s in SUSE Manager: %d\n", g, minion_id)
		} else {
			logger.Infof("Minion %s will not be added to the workflow db.\n", g)
		}
	}

	if len(minion_list) > 0 {
		/* for _, minion := range minion_list {
			result := db.FirstOrCreate(&minion, minion)
			if result.RowsAffected > 0 {
				logger.Infof("Created minion %s - %d\n", minion.Minion_Name, result.RowsAffected)
			} else {
				logger.Infof("Minion %s already exists\n", minion.Minion_Name)
			}
		} */

		returned_minions := pkg_updates.Detect_Online_Minions(SessionKey, minion_list, groupsdata)

		for _, minion_data := range returned_minions {
			//logger.Infof("------- Online check returned %s - %s\n", minion_data.Minion_Name, minion_data.Minion_Status)
			if minion_data.Minion_ID != 0 && minion_data.Minion_Status == "Online" {
				//db.Model(&minion_data).Select("Minion_Name").Updates(map[string]interface{}{"Minion_Status": minion_data.Minion_Remarks, "Minion_Remarks": "",
				//	"JobID": 0, "JobStatus": "", "Migration_Stage": "", "Migration_Stage_Status": ""})
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("JobStatus", minion_data.JobStatus)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("JobID", minion_data.JobID)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Migration_Stage", minion_data.Migration_Stage)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Remarks", minion_data.Minion_Remarks)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Migration_Stage_Status", minion_data.Migration_Stage_Status)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Status", minion_data.Minion_Status)
			}
			if minion_data.Minion_ID != 0 && minion_data.Minion_Status == "Offline" {
				minion_data.Minion_Remarks = "Offline"
				//db.Model(&minion_data).Select("Minion_Name").Updates(map[string]interface{}{"Minion_Status": minion_data.Minion_Remarks, "Minion_Remarks": "",
				//	"JobID": 0, "JobStatus": "", "Migration_Stage": "", "Migration_Stage_Status": ""})
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("JobStatus", minion_data.JobStatus)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("JobID", minion_data.JobID)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Migration_Stage", minion_data.Migration_Stage)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Remarks", minion_data.Minion_Remarks)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Migration_Stage_Status", minion_data.Migration_Stage_Status)
				db.Model(&minion_data).Where("Minion_Name = ?", minion_data.Minion_Name).Update("Minion_Status", minion_data.Minion_Status)
			}
		}

		pkg_updates.Salt_Refresh_Grains_New_by_List(SessionKey, groupsdata, returned_minions, db)
		pkg_updates.Salt_No_Upgrade_Exception_Check_New_by_List(SessionKey, groupsdata, returned_minions, db)
		pkg_updates.Salt_Disk_Space_Check_New_by_List(SessionKey, groupsdata, returned_minions, db)
		pkg_updates.Salt_Run_state_apply_by_List(groupsdata, minion_list, "pre", db)
	}
}
