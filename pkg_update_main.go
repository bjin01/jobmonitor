package main

import (
	"context"
	"os"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/pkg_updates"
	"github.com/bjin01/jobmonitor/request"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/* func ifGroupsExist(db *gorm.DB, group string) bool {
	g := new(pkg_updates.Group)
	db.Where("Group_Name = ?", group).First(&g)
	if g.Group_Name == group {
		logger.Infof("Group %s exists since %v\n", group, g.CreatedAt)
		return true
	}
	return false
} */

func GetAll_Groups(db *gorm.DB) ([]pkg_updates.Group, error) {
	var grp []pkg_updates.Group
	err := db.Preload(clause.Associations).Find(&grp).Error
	//err := db.Model(&grp).Preload("Posts").Find(&grp).Error
	return grp, err
}

func GetAll_Minions_From_DB(db *gorm.DB) ([]pkg_updates.Minion_Data, error) {
	var minion_data []pkg_updates.Minion_Data
	err := db.Preload(clause.Associations).Find(&minion_data).Error
	//err := db.Model(&grp).Preload("Posts").Find(&grp).Error
	return minion_data, err
}

func Pkg_update_groups_lookup(ctx context.Context, SUMAConfig *SUMAConfig, groupsdata *pkg_updates.Update_Groups,
	email_template_dir *email.Templates_Dir, health *bool) {

	if groupsdata.Log_Level == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	if health != nil {
		if !*health {
			logger.WithFields(logrus.Fields{
				"goroutine": ctx.Value("goroutine"),
				"username":  groupsdata.T7User,
			}).Infof("Health check failed. Skipping groups lookup.")
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
			sumaconf.Email_to = nil // or use []string{} if preferred
		}
	}
	SessionKey := new(auth.SumaSessionKey)
	var err error
	MysumaLogin := auth.Sumalogin{Login: sumaconf.User, Passwd: sumaconf.Password}
	request.Sumahost = &sumaconf.Server
	*SessionKey, err = auth.Login("auth.login", MysumaLogin)
	if err != nil {
		logger.Println(err)
		return
	}
	//email_template_directory_string := fmt.Sprintf("%s", email_template_dir.Dir)

	logger.WithFields(logrus.Fields{
		"goroutine": ctx.Value("goroutine"),
		"username":  groupsdata.T7User,
	}).Infof("Use sqlite database: %s\n", groupsdata.Sqlite_db)
	db, err := gorm.Open(gorm.Dialector(&sqlite.Dialector{DSN: groupsdata.Sqlite_db}),
		&gorm.Config{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"goroutine": ctx.Value("goroutine"),
			"username":  groupsdata.T7User,
		}).Errorf("failed to connect database")
		return
	}

	// Create the DB schema
	db.AutoMigrate(&pkg_updates.Workflow_Step{})
	db.AutoMigrate(&pkg_updates.Jobchecker_Email{})
	db.AutoMigrate(&pkg_updates.Group{})
	db.AutoMigrate(&pkg_updates.OptionalChannels{})
	db.AutoMigrate(&pkg_updates.Minion_Data{})

	//var workflow_steps []pkg_updates.Workflow_Step
	for _, g := range groupsdata.Workflow {
		new_workflow := new(pkg_updates.Workflow_Step)
		for name, step := range g {
			new_workflow.Name = name
			new_workflow.Order = step

			//workflow_steps = append(workflow_steps, *new_workflow)
			//fmt.Printf("-----------new_workflow: %+v\n", new_workflow)
			result := db.FirstOrCreate(&new_workflow, new_workflow)
			if result.RowsAffected > 0 {
				logger.WithFields(logrus.Fields{
					"goroutine": ctx.Value("goroutine"),
					"username":  groupsdata.T7User,
				}).Infof("Created workflow step %s - %d\n", name, result.RowsAffected)
			} else {
				db.Model(&new_workflow).Where("Name = ?", name).Update("Order", step)
				logger.WithFields(logrus.Fields{
					"goroutine": ctx.Value("goroutine"),
					"username":  groupsdata.T7User,
				}).Infof("Workflow Step %s already exists\n", name)
			}
		}
	}

	// Get CLM projects and environments into DB
	pkg_updates.Get_Clm_Data(SessionKey, groupsdata, db)

	/* wf, err := pkg_updates.Get_Workflow_From_DB(db)
	if err != nil {
		logger.Errorf("failed to get workflow from DB.")
		return
	}

	for _, w := range wf {
		logger.Infof("Workflow Step: %s - %d\n", w.Name, w.Order)
	} */

	for _, g := range groupsdata.Groups {

		new_group := new(pkg_updates.Group)
		new_group.Group_Name = g
		new_group.T7User = groupsdata.T7User
		new_group.Ctx_ID = groupsdata.Ctx_ID

		for _, email := range groupsdata.JobcheckerEmails {
			var jobchecker_email pkg_updates.Jobchecker_Email
			jobchecker_email.Email = email
			new_group.Email = append(new_group.Email, jobchecker_email)

		}

		//fmt.Printf("-----------new_group: %+v\n", new_group)
		result := db.FirstOrCreate(&new_group, new_group)
		if result.RowsAffected > 0 {
			logger.WithFields(logrus.Fields{
				"goroutine": ctx.Value("goroutine"),
				"username":  groupsdata.T7User,
			}).Infof("Created group %s - %d\n", g, result.RowsAffected)
		} else {
			logger.WithFields(logrus.Fields{
				"goroutine": ctx.Value("goroutine"),
				"username":  groupsdata.T7User,
			}).Infof("Group %s already exists\n", g)
			db.Model(&new_group).Where("Group_Name = ?", g).Update("T7User", groupsdata.T7User)
			db.Model(&new_group).Where("Group_Name = ?", g).Update("Ctx_ID", groupsdata.Ctx_ID)
		}

		db.Model(&new_group).Association("Email").Replace(&new_group.Email)

	}

	//db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&new_group)

	// Read
	/* mygroups, err := GetAll_Groups(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	} */

	//logger.Infof("mygroups: %+v\n", mygroups)
	/* for _, g := range mygroups {
		var email []pkg_updates.Jobchecker_Email
		db.Find(&email)
		fmt.Printf("-----------emails: %+v\n", email)
		logger.Infof("mygroup: %s\n", g.Group_Name)
		for _, e := range g.Email {
			logger.Infof("email: %s\n", e.Email)
		}
	} */

	pkg_updates.Get_Minions(SessionKey, groupsdata, db)
	pkg_updates.Salt_Refresh_Grains_New(SessionKey, groupsdata, db)

	// --------debugging start
	/* all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	post_minion_list := []pkg_updates.Minion_Data{}
	for _, g := range all_minions {
		if g.Minion_Status == "Online" {
			post_minion_list = append(post_minion_list, g)
		}
	}
	pkg_updates.Post_Migration_Debug(db, groupsdata, post_minion_list, "post_migration")
	return //for testing */
	// -----------debugging end

	pkg_updates.Salt_No_Upgrade_Exception_Check_New(SessionKey, groupsdata, db)
	pkg_updates.Salt_Disk_Space_Check_New(SessionKey, groupsdata, db)
	//logger.Debugf("---- out of if email templates dir is: %s\n", email_template_dir.Dir)

	if groupsdata.Qualifying_only {
		//set deadline to 60 seconds to allow one email sent to admins
		deadline_qualifying := time.Now().Add(time.Duration(60) * time.Second)
		//logger.Debugf("----email templates dir is: %s\n", email_template_dir.Dir)
		pkg_updates.Send_Email(ctx, groupsdata, email_template_dir, db, health, &deadline_qualifying)
		logger.WithFields(logrus.Fields{
			"goroutine": ctx.Value("goroutine"),
			"username":  groupsdata.T7User,
		}).Infof("Qualifying only is set to true. Stop the workflow here.\n")
		return
	}

	pkg_updates.Salt_Run_state_apply(groupsdata, "pre", db)

	deadline := new(time.Time)
	if groupsdata.JobcheckerTimeout == 0 {
		*deadline = time.Now().Add(time.Duration(60) * time.Minute)
	} else {
		*deadline = time.Now().Add(time.Duration(groupsdata.JobcheckerTimeout) * time.Minute)
	}

	go pkg_updates.Check_Jobs(ctx, groupsdata, SessionKey, health, db, deadline) // deadline 10min
	go pkg_updates.Start_Workflow(ctx, SessionKey, groupsdata, db, health, deadline)
	go pkg_updates.Send_Email(ctx, groupsdata, email_template_dir, db, health, deadline)

	/* all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}
	for _, g := range all_minions {
		logger.Infof("Minion in DB: %d - %s - Minion Status %s\n", g.Minion_ID, g.Minion_Name, g.Minion_Status)
		logger.Infof("Minion in DB: Ident: %s - Base Channel %s\n", g.Target_Ident, g.Target_base_channel)
		logger.Infof("Minion in DB: Remarks: %s\n", g.Minion_Remarks)
		logger.Infof("Minion in DB: Optional Channels: %v\n", g.Target_Optional_Channels)
	}
	*/
}

func Pkg_update_groups_lookup_from_file(filename string) []pkg_updates.Minion_Data {
	//get db data from filename
	logger.Infof("Use sqlite database: %s\n", filename)
	db, err := gorm.Open(gorm.Dialector(&sqlite.Dialector{DSN: filename}),
		&gorm.Config{})
	if err != nil {
		logger.Errorf("failed to initiate database connection.")
		return nil
	}

	all_minions, err := GetAll_Minions_From_DB(db)
	if err != nil {
		logger.Errorf("failed to connect database in GetAll_Minions_From_DB.")
		return nil
	}

	return all_minions
}
