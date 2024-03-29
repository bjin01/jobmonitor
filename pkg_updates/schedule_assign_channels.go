package pkg_updates

import (
	"fmt"
	"strings"
	"time"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

type Get_Channels_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
	Sid        int    `xmlrpc:"sid"`
}

/* type Get_Channels_Response struct {
	Result []struct {
		Id                 int       `xmlrpc:"id,omitempty"`
		Name               string    `xmlrpc:"name,omitempty"`
		Label              string    `xmlrpc:"label,omitempty"`
		ArchName           string    `xmlrpc:"arch_name,omitempty"`
		ArchLabel          string    `xmlrpc:"arch_label,omitempty"`
		Summary            string    `xmlrpc:"summary,omitempty"`
		Description        string    `xmlrpc:"description,omitempty"`
		ChecksumLabel      string    `xmlrpc:"checksum_label,omitempty"`
		LastModified       time.Time `xmlrpc:"last_modified,omitempty"`
		MaintainerName     string    `xmlrpc:"maintainer_name,omitempty"`
		MaintainerEmail    string    `xmlrpc:"maintainer_email,omitempty"`
		MaintainerPhone    string    `xmlrpc:"maintainer_phone,omitempty"`
		SupportPolicy      string    `xmlrpc:"support_policy,omitempty"`
		GPGKeyURL          string    `xmlrpc:"gpg_key_url,omitempty"`
		GPGKeyID           string    `xmlrpc:"gpg_key_id,omitempty"`
		GPGKeyFP           string    `xmlrpc:"gpg_key_fp,omitempty"`
		YumrepoLastSync    time.Time `xmlrpc:"yumrepo_last_sync,omitempty"`
		EndOfLife          string    `xmlrpc:"end_of_life,omitempty"`
		ParentChannelLabel string    `xmlrpc:"parent_channel_label,omitempty"`
		CloneOriginal      string    `xmlrpc:"clone_original,omitempty"`
		//ContentSources     []ContentSource `xmlrpc:"contentSources,omitempty"`
	}
} */

type Get_Channels_Response struct {
	Result []struct {
		Id                   int             `xmlrpc:"id,omitempty"`
		Name                 string          `xmlrpc:"name,omitempty"`
		Label                string          `xmlrpc:"label,omitempty"`
		Arch_name            string          `xmlrpc:"arch_name,omitempty"`
		Arch_label           string          `xmlrpc:"arch_label,omitempty"`
		Summary              string          `xmlrpc:"summary,omitempty"`
		Description          string          `xmlrpc:"description,omitempty"`
		Checksum_label       string          `xmlrpc:"checksum_label,omitempty"`
		Last_modified        time.Time       `xmlrpc:"last_modified,omitempty"`
		Maintainer_name      string          `xmlrpc:"maintainer_name,omitempty"`
		Maintainer_email     string          `xmlrpc:"maintainer_email,omitempty"`
		Maintainer_phone     string          `xmlrpc:"maintainer_phone,omitempty"`
		Support_policy       string          `xmlrpc:"support_policy,omitempty"`
		Gpg_key_url          string          `xmlrpc:"gpg_key_url,omitempty"`
		Gpg_key_id           string          `xmlrpc:"gpg_key_id,omitempty"`
		Gpg_key_fp           string          `xmlrpc:"gpg_key_fp,omitempty"`
		Yumrepo_last_sync    time.Time       `xmlrpc:"yumrepo_last_sync,omitempty"`
		End_of_life          string          `xmlrpc:"end_of_life,omitempty"`
		Parent_channel_label string          `xmlrpc:"parent_channel_label,omitempty"`
		Clone_original       string          `xmlrpc:"clone_original,omitempty"`
		ContentSources       []ContentSource `xmlrpc:"contentSources,omitempty"`
	}
}

type ContentSource struct {
	Id        int    `xmlrpc:"id,omitempty"`
	Label     string `xmlrpc:"label,omitempty"`
	SourceUrl string `xmlrpc:"sourceUrl,omitempty"`
	Type      string `xmlrpc:"type,omitempty"`
}

type Change_Channels_Request struct {
	Sessionkey         string    `xmlrpc:"sessionKey"`
	Sid                int       `xmlrpc:"sid"`
	BaseChannelLabel   string    `xmlrpc:"baseChannelLabel"`
	ChildLabels        []string  `xmlrpc:"childLabels"`
	EarliestOccurrence time.Time `xmlrpc:"earliestOccurrence"`
}

func Assign_Channels(sessionkey *auth.SumaSessionKey, groupsdata *Update_Groups, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	method := "system.listSubscribedChildChannels"

	// get Clm from db
	clm := new([]ContentLifecycleManagement)
	clm_result := db.Find(&clm)
	if clm_result.RowsAffected == 0 {
		logger.Errorf("failed to get ContentLifecycleManagement from database\n")
	}

	for i, minion := range minion_list {
		if stage != Find_Next_Stage(wf, minion) {
			continue
		} else {
			logger.Infof("Minion %s starts %s stage.\n", minion.Minion_Name, stage)
		}

		params := Get_Channels_Request{
			Sessionkey: sessionkey.Sessionkey,
			Sid:        minion.Minion_ID,
		}

		buf, err := gorillaxml.EncodeClientRequest(method, &params)
		if err != nil {
			logger.Fatalf("Encoding error: %s\n", err)
		}
		//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			logger.Fatalf("Encoding error: %s\n", err)
		}

		/* responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Fatalf("ReadAll error: %s\n", err)
		}
		logger.Infof("responseBody: %s\n", responseBody) */

		reply := new(Get_Channels_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			logger.Fatalf("Decode Get Channel Reponse body failed: %s\n", err)
		}
		//logger.Infof("Channel list %v\n", reply.Result)
		var set_channels_request Change_Channels_Request
		set_channels_request.Sessionkey = sessionkey.Sessionkey
		set_channels_request.Sid = minion.Minion_ID
		set_channels_request.EarliestOccurrence = time.Now()

		var old_base_channel_label string
		proj_label := ""
		env_label := ""
		original_base_channel_label := ""
		for _, channel := range reply.Result {
			//logger.Infof("Channel: %v\n", channel)
			var temp_base_channel_label string
			var temp_Child_label string
			update_channel_prefix := new(string)
			//logger.Infoln("Channel: ", channel.Label)
			if strings.TrimSpace(channel.Clone_original) != "" {

				if strings.TrimSpace(channel.Parent_channel_label) == "" {
					logger.Infof("Channel %s has no parent channel label\n", channel.Label)
					break
				} else {
					if original_base_channel_label == "" {
						proj_label, env_label, original_base_channel_label = Match_Project_Environment_Label(*clm, strings.TrimSpace(channel.Parent_channel_label))
					}
					logger.Debugf("matched project label: %s env label: %s, original base channel label %s\n", proj_label, env_label, original_base_channel_label)
					if proj_label == "" || env_label == "" {
						logger.Infof("Channel %s has no matching project label or environment label\n", channel.Label)
						break
					}
					//parts := strings.Split(strings.TrimSpace(channel.Parent_channel_label), "-")
					if len(groupsdata.Target_Products) > 0 {
						for _, v := range groupsdata.Target_Products {
							if len(v.Product.OptionalChildChannels) > 0 {
								var opt_channel OptionalChannels
								for _, optchannel := range v.Product.OptionalChildChannels {
									if strings.Contains(strings.TrimSpace(channel.Label), strings.TrimSpace(optchannel.Old_Channel)) {
										if strings.TrimSpace(optchannel.New_Channel) != "" {
											if v.Product.Clm_Project_Label != "" {
												new_opt_channel_label := fmt.Sprintf("%s-%s-%s",
													v.Product.Clm_Project_Label, env_label, strings.TrimSpace(optchannel.New_Channel))

												opt_channel.Channel_Label = new_opt_channel_label
												result := db.FirstOrCreate(&opt_channel)
												if result.RowsAffected > 0 {
													logger.Debugf("Optional channel %s is created\n", new_opt_channel_label)

												} else {
													db.Model(&minion).Association("Target_Optional_Channels").Unscoped().Clear()
												}
												db.Model(&minion_list[i]).Association("Target_Optional_Channels").Append(&opt_channel)
												logger.Debugf("Optional channel %s is assigned to %s\n", new_opt_channel_label, minion.Minion_Name)
											} else {
												new_opt_channel_label := strings.TrimSpace(optchannel.New_Channel)
												opt_channel.Channel_Label = new_opt_channel_label
												result := db.FirstOrCreate(&opt_channel)
												if result.RowsAffected > 0 {
													logger.Debugf("Optional channel %s is created\n", new_opt_channel_label)

												} else {
													db.Model(&minion).Association("Target_Optional_Channels").Unscoped().Clear()
												}
												db.Model(&minion_list[i]).Association("Target_Optional_Channels").Append(&opt_channel)
												logger.Debugf("Optional channel %s is assigned to %s\n", new_opt_channel_label, minion.Minion_Name)
											}
										}
									}
								}
							}
						}
					}

					if env_label != "" {
						db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Clm_Stage", env_label)
						logger.Infof("Minion %s is at content lifecycle management stage: %s-%s\n", minion.Minion_Name, proj_label, env_label)

					} else {
						logger.Infof("%s: Channel %s could not be parsed.\n", minion.Minion_Name, channel.Label)
						db.Model(&minion).Where("Minion_Name = ?", minion.Minion_Name).Update("Clm_Stage", "")
					}
				}

				for _, group := range groupsdata.Assigne_channels {
					if strings.TrimSpace(group.Assigne_Channel.Current_base_channel) == strings.TrimSpace(channel.Parent_channel_label) {
						old_base_channel_label = channel.Parent_channel_label
						if strings.TrimSpace(group.Assigne_Channel.New_base_prefix) != "" {
							update_channel_prefix = &group.Assigne_Channel.New_base_prefix
							temp_base_channel_label = fmt.Sprintf("%s%s", *update_channel_prefix, original_base_channel_label)
							_, _, original_child_channel_label := Match_Project_Environment_Label(*clm, strings.TrimSpace(channel.Label))
							temp_Child_label = fmt.Sprintf("%s%s", *update_channel_prefix, original_child_channel_label)
						} else {
							*update_channel_prefix = ""
							temp_base_channel_label = original_base_channel_label
							_, _, original_child_channel_label := Match_Project_Environment_Label(*clm, strings.TrimSpace(channel.Label))
							temp_Child_label = original_child_channel_label
						}
					}
				}

				set_channels_request.BaseChannelLabel = temp_base_channel_label
				set_channels_request.ChildLabels = append(set_channels_request.ChildLabels, temp_Child_label)

			} else {
				if len(groupsdata.Target_Products) > 0 {
					for _, v := range groupsdata.Target_Products {
						if len(v.Product.OptionalChildChannels) > 0 {
							var opt_channel OptionalChannels
							for _, optchannel := range v.Product.OptionalChildChannels {
								if strings.Contains(strings.TrimSpace(channel.Label), strings.TrimSpace(optchannel.Old_Channel)) {
									if strings.TrimSpace(optchannel.New_Channel) != "" {

										new_opt_channel_label := strings.TrimSpace(optchannel.New_Channel)
										opt_channel.Channel_Label = new_opt_channel_label
										result := db.FirstOrCreate(&opt_channel)
										if result.RowsAffected > 0 {
											logger.Debugf("Optional channel %s is created\n", new_opt_channel_label)

										} else {
											db.Model(&minion).Association("Target_Optional_Channels").Unscoped().Clear()
										}
										db.Model(&minion_list[i]).Association("Target_Optional_Channels").Append(&opt_channel)
										db.Model(&minion_list[i]).Where("Minion_Name = ?", minion.Minion_Name).Update("Clm_Stage", "")
										logger.Debugf("Optional channel %s is assigned to %s\n", new_opt_channel_label, minion.Minion_Name)

									}
								}
							}
						}
					}
				}

				for _, group := range groupsdata.Assigne_channels {

					if strings.Contains(strings.TrimSpace(group.Assigne_Channel.Current_base_channel),
						strings.TrimSpace(channel.Parent_channel_label)) {
						old_base_channel_label = channel.Parent_channel_label
						if strings.TrimSpace(group.Assigne_Channel.New_base_prefix) != "" {
							update_channel_prefix = &group.Assigne_Channel.New_base_prefix
							temp_base_channel_label = fmt.Sprintf("%s%s", *update_channel_prefix, channel.Parent_channel_label)
							temp_Child_label = fmt.Sprintf("%s%s", *update_channel_prefix, channel.Label)
						} else {
							temp_base_channel_label = channel.Parent_channel_label
							temp_Child_label = channel.Label
						}
						set_channels_request.BaseChannelLabel = temp_base_channel_label
						set_channels_request.ChildLabels = append(set_channels_request.ChildLabels, temp_Child_label)
					}
				}
			}

			/* logger.Infof("Channel ID: %v\n", channel.Id)
			logger.Infof("Channel Name: %v\n", channel.Name)
			logger.Infof("Channel Label: %v\n", channel.Label)
			logger.Infof("Channel Parent Label: %v\n", channel.Parent_channel_label)
			logger.Infof("Channel Clone_original Name: %v\n", channel.Clone_original) */
			//logger.Infoln()
		}

		if strings.TrimSpace(set_channels_request.BaseChannelLabel) != "" {
			if old_base_channel_label == set_channels_request.BaseChannelLabel {
				logger.Debugf("Existing %s is already assigned on %s\n", set_channels_request.BaseChannelLabel,
					minion.Minion_Name)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
				continue
			}

			//logger.Infof("Assigne %s including child channels for: %s\n",
			//	set_channels_request.BaseChannelLabel, minion.Minion_Name)
			buf, err := gorillaxml.EncodeClientRequest("system.scheduleChangeChannels", &set_channels_request)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
			resp, err := request.MakeRequest(buf)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}

			/* responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatalf("ReadAll error: %s\n", err)
			}
			logger.Infof("responseBody: %s\n", responseBody) */

			reply := new(Generic_Job_Response)
			err = gorillaxml.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				logger.Fatalf("Decode scheduleChangeChannels Job response body failed: %s\n", err)
			}
			logger.Debugf("scheduleChangeChannels JobID: %d\n", reply.JobID)
			var host_info Host_Job_Info
			host_info.Assigne_Channels_Job.JobID = reply.JobID
			host_info.Assigne_Channels_Job.JobStatus = "Scheduled"

			if reply.JobID > 0 {
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobID", reply.JobID)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("JobStatus", "pending")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "Scheduled")
				db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Target_base_channel", set_channels_request.BaseChannelLabel)
			}

		} else {
			logger.Debugf("System is already on original channels. %s\n", minion.Minion_Name)
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage", stage)
			db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Migration_Stage_Status", "completed")
		}
	}

}
