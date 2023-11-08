package pkg_updates

import (
	"log"
	"strings"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/request"
	gorillaxml "github.com/divan/gorilla-xmlrpc/xml"
)

//This func is for finding the migration target for each minion at the beginning of the migration process.
//If no proper migration target is found the minion will be excluded from further processing.
func Find_MigrationTarget_New(sessionkey *auth.SumaSessionKey, minionid int, UserData *Update_Groups) (ident string, migrate_base_channel string) {
	method := "system.listMigrationTargets"
	ident = ""
	migrate_base_channel = ""

	var params ListMigrationTarget_Request
	params.Sessionkey = sessionkey.Sessionkey
	params.Sid = minionid
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
	defer resp.Body.Close()
	reply := new(ListMigrationTarget_Response)
	err = gorillaxml.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Default().Printf("Decode ListMigrationTarget response body failed: %s\n", err)
	}
	for _, target := range reply.Result {
		//split_result := Convert_String_to_maps(target.Friendly)
		if UserData.Target_Products != nil {

			for _, v := range UserData.Target_Products {
				//logger.Infof("v: %s\n", v)
				//logger.Infof("target ident: %s vs. defined ident: %s\n", target.Ident, v.Product.Ident)
				//logger.Infof("target product: %s vs. defined product: %s\n", split_result["base"], v.Product.Name)
				if strings.Contains(target.Ident, v.Product.Ident) {
					logger.Debugf("Ident Match found %s\n", target.Ident)
					ident = target.Ident
					migrate_base_channel = v.Product.Base_Channel
				}
			}
		}
	}
	return ident, migrate_base_channel
}
