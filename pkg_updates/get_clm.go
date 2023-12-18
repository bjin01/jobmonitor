package pkg_updates

import (
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

type ContentLifecycleManagement struct {
	gorm.Model
	Clm_Project_Label             string
	Clm_Project_Environment_Label string
}

type List_Projects_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
}

type List_Projects_Response struct {
	Result []struct {
		Id               int       `xmlrpc:"id"`
		Label            string    `xmlrpc:"label"`
		Name             string    `xmlrpc:"name"`
		Description      string    `xmlrpc:"description"`
		LastBuildDate    time.Time `xmlrpc:"lastBuildDate"`
		OrgId            int       `xmlrpc:"orgId"`
		FirstEnvironment string    `xmlrpc:"firstEnvironment"`
	}
}

type List_Project_Environments_Request struct {
	Sessionkey   string `xmlrpc:"sessionKey"`
	ProjectLabel string `xmlrpc:"projectLabel"`
}

type List_Project_Environments_Response struct {
	Result []struct {
		Id                       int       `xmlrpc:"id"`
		Label                    string    `xmlrpc:"label"`
		Name                     string    `xmlrpc:"name"`
		Description              string    `xmlrpc:"description"`
		Status                   string    `xmlrpc:"status"`
		LastBuildDate            time.Time `xmlrpc:"lastBuildDate"`
		ContentProjectLabel      string    `xmlrpc:"contentProjectLabel"`
		PreviousEnvironmentLabel string    `xmlrpc:"previousEnvironmentLabel"`
		NextEnvironmentLabel     string    `xmlrpc:"nextEnvironmentLabel"`
	}
}

func Get_Clm_Data(sessionkey *auth.SumaSessionKey, UserData *Update_Groups, db *gorm.DB) {

	err := db.AutoMigrate(&ContentLifecycleManagement{})
	if err != nil {
		logger.Errorf("failed to create ContentLifecycleManagement table in DB: %s\n", err)
		return
	}

	method := "contentmanagement.listProjects"

	list_projects_request := List_Projects_Request{
		Sessionkey: sessionkey.Sessionkey,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &list_projects_request)
	if err != nil {
		logger.Infof("Encoding list_projects_request error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Infof("Encoding list_projects_request error: %s\n", err)
	}

	reply := new(List_Projects_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Infof("Decode list_projects_request response body failed: %s\n", err)
	}

	if len(reply.Result) > 0 {
		for _, project := range reply.Result {
			List_Project_Environments(sessionkey, db, project.Label)
		}
	} else {
		logger.Infof("No CLM projects found\n")
	}

	//result := db.FirstOrCreate(&new_workflow, new_workflow)
}

func List_Project_Environments(sessionkey *auth.SumaSessionKey, db *gorm.DB, proj_label string) {
	method := "contentmanagement.listProjectEnvironments"

	list_project_environments_request := List_Project_Environments_Request{
		Sessionkey:   sessionkey.Sessionkey,
		ProjectLabel: proj_label,
	}

	buf, err := gorillaxml.EncodeClientRequest(method, &list_project_environments_request)
	if err != nil {
		logger.Infof("Encoding list_project_environments_request error: %s\n", err)
	}

	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Infof("Encoding list_project_environments_request error: %s\n", err)
	}

	reply := new(List_Project_Environments_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Infof("Decode List_Project_Environments_Response response body failed: %s\n", err)
	}

	if len(reply.Result) > 0 {
		logger.Debugf("Project %s has %d environments\n", proj_label, len(reply.Result))

		for _, environment := range reply.Result {
			if environment.ContentProjectLabel != "" && environment.Label != "" {
				new_clm := ContentLifecycleManagement{
					Clm_Project_Label:             proj_label,
					Clm_Project_Environment_Label: environment.Label,
				}
				result := db.FirstOrCreate(&new_clm, new_clm)
				if result.RowsAffected > 0 {
					logger.Infof("Created CLM Project %s - %d\n", proj_label, result.RowsAffected)
				} else {
					db.Model(&new_clm).Where("Clm_Project_Label = ?", proj_label).Update("Clm_Project_Environment_Label", environment.Label)
					logger.Infof("CLM Project %s already exists\n", proj_label)
				}
			}
		}
	}

	return

}

func Match_Project_Environment_Label(all_projects []ContentLifecycleManagement, channel_label string) (project_label, longest_env_label, original_channel_label string) {

	longest_env_label = ""
	project_label = ""
	original_channel_label = channel_label

	if len(all_projects) > 0 {
		temp_result := make(map[string][]string)
		for _, project := range all_projects {
			proj_env_label := project.Clm_Project_Label + "-" + project.Clm_Project_Environment_Label + "-"
			if strings.HasPrefix(channel_label, proj_env_label) {
				temp_result[project.Clm_Project_Label] = append(temp_result[project.Clm_Project_Label], project.Clm_Project_Environment_Label)
			}
		}

		// need below to get the longest env label because environment label can be like "dev", "dev1", "dev2", etc.
		if len(temp_result) > 0 {
			for project, envs := range temp_result {
				if len(envs) > 0 {
					for _, env := range envs {
						if len(env) > len(longest_env_label) {
							longest_env_label = env
							project_label = project
						}
					}
				}

			}

			substring := project_label + "-" + longest_env_label + "-"
			index := strings.Index(channel_label, substring)
			if index != -1 {
				original_channel_label = channel_label[index+len(substring):]
			}

		}
	}
	return project_label, longest_env_label, original_channel_label
}
