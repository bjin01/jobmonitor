package saltapi

type Login_Response struct {
	Return []struct {
		Token  string   `json:"token"`
		Expire float64  `json:"expire"`
		Start  float64  `json:"start"`
		User   string   `json:"user"`
		Eauth  string   `json:"eauth"`
		Perms  []string `json:"perms"`
	} `json:"return"`
}

type Salt_Data struct {
	Username                        string   `json:"username"`
	Password                        string   `json:"password"`
	SaltMaster                      string   `json:"salt_master"`
	SaltApi_Port                    int      `json:"salt_api_port"`
	Token                           string   `json:"token,omitempty"`
	Salt_Client_Type                string   `json:"salt_client_type,omitempty"`
	SaltCmd                         string   `json:"salt_cmd,omitempty"`
	Salt_diskspace_grains_key       string   `json:"salt_diskspace_grains_key,omitempty"`
	Salt_diskspace_grains_value     string   `json:"salt_diskspace_grains_value,omitempty"`
	Salt_no_upgrade_exception_key   string   `json:"salt_no_upgrade_exception_key"`
	Salt_no_upgrade_exception_value string   `json:"salt_no_upgrade_exception_value"`
	Target_List                     []string `json:"target_list,omitempty"`
	Arg                             []string `json:"arg,omitempty"`
	Return                          []byte   `json:"return,omitempty"`
	Online_Minions                  []string `json:"online_minions,omitempty"`
	Offline_Minions                 []string `json:"offline_minions,omitempty"`
	Jid                             string   `json:"jid,omitempty"`
}

type SaltJob_Data struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SaltMaster   string `json:"salt_master"`
	SaltApi_Port int    `json:"salt_api_port"`
	Token        string `json:"token,omitempty"`
	Jid          string `json:"jid,omitempty"`
}

type Salt_Request struct {
	Client   string   `json:"client"`
	Tgt      []string `json:"tgt"`
	Tgt_type string   `json:"tgt_type"`
	Fun      string   `json:"fun"`
	Arg      []string `json:"arg"`
}

type Salt_Request_Async struct {
	Tgt      []string `json:"tgt"`
	Tgt_type string   `json:"tgt_type"`
	Fun      string   `json:"fun"`
	Arg      []string `json:"arg"`
}

//{"return": [{"up": ["pxesap01.bo2go.home", "pxesap02.bo2go.home"], "down": ["jupiter.bo2go.home"]}]}
type Runner_Manage_Status_Response struct {
	Return []struct {
		Down []string `json:"down"`
		Up   []string `json:"up"`
	} `json:"return"`
}

type Salt_Async_Response struct {
	Return []struct {
		Jid     string   `json:"jid"`
		Minions []string `json:"minions"`
	} `json:"return"`
}

type ResultData struct {
	Return  map[string]Return_Inner `json:"return"`
	RetCode int                     `json:"retcode"`
	Success bool                    `json:"success"`
	Out     string                  `json:"out"`
}

type Return_Inner struct {
	Name    interface{} `json:"name"`
	Changes interface{} `json:"changes"`
	Result  bool        `json:"result"`
	Comment string      `json:"comment"`
	Sls     string      `json:"__sls__"`
	Sls_Id  string      `json:"__id__"`
}

type SaltResponse struct {
	Return []SaltJob `json:"return"`
}

type SaltJob struct {
	JID        string                 `json:"jid"`
	Function   string                 `json:"Function"`
	Arguments  []string               `json:"Arguments"`
	Target     []string               `json:"Target"`
	TargetType string                 `json:"Target-type"`
	User       string                 `json:"User"`
	Minions    []string               `json:"Minions"`
	StartTime  string                 `json:"StartTime"`
	Result     map[string]interface{} `json:"Result"`
}
