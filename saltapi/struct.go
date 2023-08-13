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
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	SaltMaster       string   `json:"salt_master"`
	SaltApi_Port     int      `json:"salt_api_port"`
	Token            string   `json:"token,omitempty"`
	Salt_Client_Type string   `json:"salt_client_type,omitempty"`
	SaltCmd          string   `json:"salt_cmd,omitempty"`
	Target_List      []string `json:"target_list,omitempty"`
	Arg              []string `json:"arg,omitempty"`
	Return           []byte   `json:"return,omitempty"`
	Online_Minions   []string `json:"online_minions,omitempty"`
	Offline_Minions  []string `json:"offline_minions,omitempty"`
}

//{"return": [{"up": ["pxesap01.bo2go.home", "pxesap02.bo2go.home"], "down": ["jupiter.bo2go.home"]}]}
type Runner_Manage_Status_Response struct {
	Return []struct {
		Down []string `json:"down"`
		Up   []string `json:"up"`
	} `json:"return"`
}
