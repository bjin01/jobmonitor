package spmigration

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

type ListMigrationTarget_Request struct {
	Sessionkey                          string `xmlrpc:"sessionKey"`
	Sid                                 int    `xmlrpc:"sid"`
	ExcludeTargetWhereMissingSuccessors bool   `xmlrpc:"excludeTargetWhereMissingSuccessors"`
}

type ListMigrationTarget_Response struct {
	Result []struct {
		Ident    string `xmlrpc:"ident,omitempty"`
		Friendly string `xmlrpc:"friendly,omitempty"`
	}
}

type ListAllChannels_Request struct {
	Sessionkey string `xmlrpc:"sessionKey"`
}

type ListAllChannels_Response struct {
	Result []struct {
		Id            int    `xmlrpc:"id,omitempty"`
		Name          string `xmlrpc:"name,omitempty"`
		Label         string `xmlrpc:"label,omitempty"`
		Arch_name     string `xmlrpc:"arch_name,omitempty"`
		Provider_name string `xmlrpc:"provider_name,omitempty"`
		Packages      int    `xmlrpc:"packages,omitempty"`
		Systems       int    `xmlrpc:"systems,omitempty"`
	}
}

func (t *Target_Minions) ListMigrationTarget(sessionkey *auth.SumaSessionKey, UserData *Migration_Groups) {
	method := "system.listMigrationTargets"
	//allchannels := List_All_Channels(sessionkey)

	for i, minion := range t.Minion_List {
		if minion.Migration_Stage_Status == "Completed" &&
			(minion.Migration_Stage == "Pkg_Refresh" || minion.Migration_Stage == "Reboot") {
			log.Printf("Minion %s is ready for spmigration dryrun.\n", minion.Minion_Name)
		} else {
			log.Printf("Minion %s is not ready for spmigration dryrun, previous step <Pkg_Refresh> not successful or finished.\n", minion.Minion_Name)
			continue
		}
		//t.Minion_List[i].Target_base_channel = UserData.Target_base_channel
		var params ListMigrationTarget_Request
		params.Sessionkey = sessionkey.Sessionkey
		params.Sid = minion.Minion_ID
		params.ExcludeTargetWhereMissingSuccessors = false
		buf, err := gorillaxml.EncodeClientRequest(method, &params)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		resp, err := request.MakeRequest(buf)
		if err != nil {
			log.Fatalf("Encoding error: %s\n", err)
		}
		//fmt.Printf("buffer: %s\n", string(buf))
		//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
		/* responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll error: %s\n", err)
		}
		fmt.Printf("responseBody: %s\n", responseBody) */
		reply := new(ListMigrationTarget_Response)
		err = gorillaxml.DecodeClientResponse(resp.Body, reply)
		if err != nil {
			log.Fatalf("Decode ListMigrationTarget response body failed: %s\n", err)
		}
		for _, target := range reply.Result {
			split_result := Convert_String_to_maps(target.Friendly)
			if UserData.Target_Products != nil {

				for _, v := range UserData.Target_Products {
					//fmt.Printf("v: %s\n", v)
					//fmt.Printf("target: %s value: %s\n", target.Ident, v.Product.Ident)
					//fmt.Printf("target: %s value: %s\n", split_result["base"], v.Product.Name)
					if strings.Contains(target.Ident, v.Product.Ident) {
						fmt.Printf("%s\n", minion.Minion_Name)
						fmt.Printf("Found matching Target product ident: %s\n", target.Ident)
						fmt.Printf("Found matching Target product name: %s\n", split_result["base"])
						fmt.Printf("Found matching Target product base channel: %s\n", v.Product.Base_Channel)
						fmt.Println()

						// we do this so that every system gets the base channel with
						// the clm project and same environment as pior service pack clm env. set.

						for _, v := range t.Minion_Environment_List {
							for hostname, environment := range v {
								//fmt.Printf("hostname: [%s] original_environment: [%s]\n", k, v)
								if hostname == minion.Minion_Name {
									minion.Host_Job_Info.Channel_Environment = environment
									log.Printf("Set hostname: [%s] to original environment label: [%s]\n", minion.Minion_Name, environment)
								}
							}
						}

						if minion.Host_Job_Info.Channel_Environment != "" &&
							v.Product.Clm_Project_Label != "" &&
							v.Product.Base_Channel != "" {
							target_base_channel := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(v.Product.Clm_Project_Label),
								minion.Host_Job_Info.Channel_Environment, strings.TrimSpace(v.Product.Base_Channel))
							t.Minion_List[i].Target_base_channel = target_base_channel
							t.Minion_List[i].Target_Ident = target.Ident
							if v.Product.OptionalChildChannels != nil {
								for _, child := range v.Product.OptionalChildChannels {
									log.Printf("%s: Add clm optional channel to schedule spmigration: %s\n",
										minion.Minion_Name, child)
									optional_channel := fmt.Sprintf("%s-%s-%s", strings.TrimSpace(v.Product.Clm_Project_Label),
										minion.Host_Job_Info.Channel_Environment, strings.TrimSpace(child))
									t.Minion_List[i].Target_Optional_Channels = append(t.Minion_List[i].Target_Optional_Channels, optional_channel)
								}
							}
						} else {
							// if the env is not provided or empty then we use the base channel only.
							t.Minion_List[i].Target_base_channel = strings.TrimSpace(v.Product.Base_Channel)
							t.Minion_List[i].Target_Ident = target.Ident
							if v.Product.OptionalChildChannels != nil {
								for _, child := range v.Product.OptionalChildChannels {
									log.Printf("%s: Add optional channel to schedule spmigration: %s\n",
										minion.Minion_Name, child)
									optional_channel := strings.TrimSpace(child)
									t.Minion_List[i].Target_Optional_Channels = append(t.Minion_List[i].Target_Optional_Channels, optional_channel)
								}
							}
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
		fmt.Printf("%s: %s\n", key, value)
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
	//fmt.Println(intSlice)
	return intSlice
}

func List_All_Channels(sessionkey *auth.SumaSessionKey) *ListAllChannels_Response {
	method := "channel.listAllChannels"
	var params ListAllChannels_Request
	params.Sessionkey = sessionkey.Sessionkey

	buf, err := gorillaxml.EncodeClientRequest(method, &params)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	resp, err := request.MakeRequest(buf)
	if err != nil {
		log.Fatalf("Encoding error: %s\n", err)
	}
	//fmt.Printf("buffer: %s\n", string(buf))
	//fmt.Printf("buffer: %s\n", fmt.Sprintf(string(buf)))
	/* responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %s\n", err)
	}
	fmt.Printf("responseBody: %s\n", responseBody) */
	reply := new(ListAllChannels_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Fatalf("Decode List_All_Channels response body failed: %s\n", err)
	}
	//fmt.Printf("List_All_Vendor_Channels: %v\n", reply.Result)
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

	fmt.Println("major_version:", majorVersion)
	fmt.Println("service_pack:", servicePack)
	fmt.Println("name_of_product:", nameOfProduct)
}
