package main

import (
	"os"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/pkg_updates"
	"github.com/bjin01/jobmonitor/request"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ifGroupsExist(db *gorm.DB, group string) bool {
	g := new(pkg_updates.Group)
	db.Where("Group_Name = ?", group).First(&g)
	if g.Group_Name == group {
		logger.Infof("Group %s exists since %v\n", group, g.CreatedAt)
		return true
	}
	return false
}

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

func Pkg_update_groups_lookup(SUMAConfig *SUMAConfig, groupsdata *pkg_updates.Update_Groups,
	email_template_dir *email.Templates_Dir, health *bool) {

	if health != nil {
		if *health == false {
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
		&gorm.Config{})
	if err != nil {
		logger.Errorf("failed to connect database")
		return
	}

	// Create the DB schema
	db.AutoMigrate(&pkg_updates.Jobchecker_Email{})
	db.AutoMigrate(&pkg_updates.Group{})
	db.AutoMigrate(&pkg_updates.OptionalChannels{})
	db.AutoMigrate(&pkg_updates.Minion_Data{})

	for _, g := range groupsdata.Groups {

		new_group := new(pkg_updates.Group)
		new_group.Group_Name = g
		new_group.T7User = groupsdata.T7User

		for _, email := range groupsdata.JobcheckerEmails {
			var jobchecker_email pkg_updates.Jobchecker_Email
			jobchecker_email.Email = email
			new_group.Email = append(new_group.Email, jobchecker_email)

		}

		//fmt.Printf("-----------new_group: %+v\n", new_group)
		result := db.FirstOrCreate(&new_group, new_group)
		if result.RowsAffected > 0 {
			logger.Infof("Created group %s - %d\n", g, result.RowsAffected)
		} else {
			logger.Infof("Group %s already exists\n", g)
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
	pkg_updates.Salt_No_Upgrade_Exception_Check_New(SessionKey, groupsdata, db)
	pkg_updates.Salt_Disk_Space_Check_New(SessionKey, groupsdata, db)
	pkg_updates.Salt_Run_state_apply(SessionKey, groupsdata, "pre", db)
	pkg_updates.Send_Email(groupsdata, email_template_dir, db)
	if groupsdata.Include_Spmigration {
		pkg_updates.Assign_Channels(SessionKey, groupsdata, db)
		//Check_Assigne_Channels_Jobs(sessionkey, health) // deadline 10min
	}

	all_minions, err := GetAll_Minions_From_DB(db)
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

}
