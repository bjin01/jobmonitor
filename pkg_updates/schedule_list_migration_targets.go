package pkg_updates

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
	"gorm.io/gorm"
)

func ListMigrationTarget(sessionkey *auth.SumaSessionKey, UserData *Update_Groups, db *gorm.DB, wf []Workflow_Step, minion_list []Minion_Data, stage string) {
	method := "system.listMigrationTargets"
	//allchannels := List_All_Channels(sessionkey)

	for _, minion := range minion_list {
		//get minion stage fromo DB
		result := db.Where(&Minion_Data{Minion_Name: minion.Minion_Name}).First(&minion)
		if result.Error != nil {
			logger.Errorf("failed to get minion %s from database\n", minion.Minion_Name)
			return
		}

		logger.Infof("Minion %s stage is %s\n", minion.Minion_Name, minion.Migration_Stage)

		if stage == Find_Next_Stage(wf, minion) {
			logger.Infof("Minion %s starts %s stage.\n", minion.Minion_Name, stage)

			//t.Minion_List[i].Target_base_channel = UserData.Target_base_channel
			var params ListMigrationTarget_Request
			params.Sessionkey = sessionkey.Sessionkey
			params.Sid = minion.Minion_ID
			params.ExcludeTargetWhereMissingSuccessors = false
			buf, err := gorillaxml.EncodeClientRequest(method, &params)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
			resp, err := request.MakeRequest(buf)
			if err != nil {
				logger.Fatalf("Encoding error: %s\n", err)
			}
			//logger.Infof("buffer: %s\n", string(buf))
			//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
			/* responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatalf("ReadAll error: %s\n", err)
			}
			logger.Infof("responseBody: %s\n", responseBody) */
			reply := new(ListMigrationTarget_Response)
			err = gorillaxml.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				logger.Fatalf("Decode ListMigrationTarget response body failed: %s\n", err)
			}
			for _, target := range reply.Result {
				split_result := Convert_String_to_maps(target.Friendly)
				if UserData.Target_Products != nil {

					for _, v := range UserData.Target_Products {
						//logger.Infof("v: %s\n", v)
						//logger.Infof("target: %s value: %s\n", target.Ident, v.Product.Ident)
						//logger.Infof("target: %s value: %s\n", split_result["base"], v.Product.Name)
						if strings.Contains(target.Ident, v.Product.Ident) {
							/* logger.Infof("%s\n", minion.Minion_Name)
							logger.Infof("Found matching Target product ident: %s\n", target.Ident)
							logger.Infof("Found matching Target product name: %s\n", split_result["base"])
							logger.Infof("Found matching Target product base channel: %s\n", v.Product.Base_Channel)
							logger.Infoln() */

							// we do this so that every system gets the base channel with
							// the clm project and same environment as pior service pack clm env. set.

							if minion.Clm_Stage != "" && v.Product.Clm_Project_Label != "" && v.Product.Base_Channel != "" {
								target_base_channel := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(v.Product.Clm_Project_Label), minion.Clm_Stage, strings.TrimSpace(v.Product.Base_Channel))
								db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Target_Ident", target.Ident)
								db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Target_base_channel", target_base_channel)

								/* if v.Product.OptionalChildChannels != nil {
									new_spmigration_optional_channels := []string{}
									for _, child := range v.Product.OptionalChildChannels {
										if len(minion.Target_Optional_Channels) != 0 {
											for _, minion_optional_channel := range minion.Target_Optional_Channels {
												if strings.Contains(minion_optional_channel.Channel_Label, child.Old_Channel) {
													formatted_new_optional_channel_label := fmt.Sprintf("%s-%s-%s",
														strings.TrimSpace(v.Product.Clm_Project_Label), minion.Clm_Stage,
														strings.TrimSpace(child.New_Channel))

													new_spmigration_optional_channels = append(new_spmigration_optional_channels, formatted_new_optional_channel_label)
													logger.Infof("%s: Add clm optional channel to schedule spmigration: %s\n",
														minion.Minion_Name, child.New_Channel)
												}
											}
										}

									}
								} */
							} else {
								// if the env is not provided or empty then we use the base channel only.
								db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Target_Ident", target.Ident)
								db.Model(&Minion_Data{}).Where("Minion_Name = ?", minion.Minion_Name).Update("Target_base_channel", v.Product.Base_Channel)

								/* if v.Product.OptionalChildChannels != nil {
									new_spmigration_optional_channels := []string{}
									for _, child := range v.Product.OptionalChildChannels {
										if len(minion.Target_Optional_Channels) != 0 {
											for _, minion_optional_channel := range minion.Target_Optional_Channels {
												if strings.Contains(minion_optional_channel.Channel_Label, child.Old_Channel) {
													new_spmigration_optional_channels = append(new_spmigration_optional_channels, child.New_Channel)
													logger.Infof("%s: Add clm optional channel to schedule spmigration: %s\n",
														minion.Minion_Name, child.New_Channel)
												}
											}
										}

									}
								} */
							}

						}
					}
				} else {
					log.Default().Printf("%s\n", minion.Minion_Name)
					log.Default().Println("No target products provided.")
					log.Default().Printf("Possible target ident: %s.", target.Ident)
					log.Default().Printf("Possible target base: %s.", split_result["base"])
				}
			}
		}
	}
}

func Convert_String_to_maps(mystring string) map[string]string {
	/* mystring := "[base: SUSE Linux Enterprise Server for SAP Applications 15 SP5 x86_64,
	addon: Desktop Applications Module 15 SP5 x86_64,
	SUSE Linux Enterprise Live Patching 15 SP5 x86_64,
	Web and Scripting Module 15 SP5 x86_64,
	Basesystem Module 15 SP5 x86_64, SAP Applications Module 15 SP5 x86_64,
	Server Applications Module 15 SP5 x86_64, SUSE Manager Client Tools for SLE 15 x86_64,
	Python 3 Module 15 SP5 x86_64,
	SUSE Linux Enterprise High Availability Extension 15 SP5 x86_64]" */

	// Remove brackets from the string
	mystring = strings.TrimPrefix(mystring, "[")
	mystring = strings.TrimSuffix(mystring, "]")

	// Split the string into key-value pairs
	pairs := strings.Split(mystring, ", ")

	// Create a map to store the key-value pairs
	m := make(map[string]string)

	// Iterate over each pair and populate the map
	for _, pair := range pairs {
		kv := strings.Split(pair, ": ")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[key] = value
		}
	}

	// Print the resulting map
	/* for key, value := range m {
		logger.Infof("%s: %s\n", key, value)
	} */
	return m
}

func Convert_String_IntSlices(mystring string) []int {

	// Remove brackets from the string
	mystring = strings.TrimPrefix(mystring, "[")
	mystring = strings.TrimSuffix(mystring, "]")

	// Split the string into individual integers
	intStrs := strings.Split(mystring, ",")

	// Create a slice to store the integers
	intSlice := make([]int, 0, len(intStrs))

	// Convert each string to an integer and append it to the slice
	for _, s := range intStrs {
		if i, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			intSlice = append(intSlice, i)
		}
	}

	// Print the resulting slice
	//logger.Infoln(intSlice)
	return intSlice
}

func List_All_Channels(sessionkey *auth.SumaSessionKey) *ListAllChannels_Response {
	method := "channel.listAllChannels"
	var params ListAllChannels_Request
	params.Sessionkey = sessionkey.Sessionkey

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		logger.Fatalf("Encoding error: %s\n", err)
	}
	//logger.Infof("buffer: %s\n", string(buf))
	//logger.Infof("buffer: %s\n", fmt.Sprintf(string(buf)))
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("ReadAll error: %s\n", err)
	}
	logger.Infof("responseBody: %s\n", responseBody) */
	reply := new(ListAllChannels_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		logger.Fatalf("Decode List_All_Channels response body failed: %s\n", err)
	}
	//logger.Infof("List_All_Vendor_Channels: %v\n", reply.Result)
	return reply
}

func Parse_Product_info(inputString string) {
	//inputString := "SUSE Linux Enterprise Server for SAP Applications 15 SP5 x86_64"

	// Regular expressions to match the major_version and service_pack
	majorVersionRegex := regexp.MustCompile(`(\d+)`)
	servicePackRegex := regexp.MustCompile(`SP(\d+)`)

	// Extracting the major_version
	majorVersion := majorVersionRegex.FindStringSubmatch(inputString)[1]

	// Extracting the service_pack
	servicePack := servicePackRegex.FindStringSubmatch(inputString)[1]

	// Extracting the name_of_product
	nameOfProduct := strings.TrimSpace(strings.Split(inputString, majorVersion)[0])

	logger.Infoln("major_version:", majorVersion)
	logger.Infoln("service_pack:", servicePack)
	logger.Infoln("name_of_product:", nameOfProduct)
}
