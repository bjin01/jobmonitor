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
	Username     string `json:"username"`
	Password     string `json:"password"`
	SaltMaster   string `json:"salt_master"`
	SaltApi_Port int    `json:"salt_api_port"`
	Token        string `json:"token,omitempty"`
}
