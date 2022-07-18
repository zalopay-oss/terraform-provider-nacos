package client

type Configuration struct {
	Namespace   string `json:"tenant"`
	Group       string `json:"group"`
	Key         string `json:"dataId"`
	Value       string `json:"content"`
	Description string `json:"desc"`
}

type ConfigurationId struct {
	Namespace string
	Group     string
	Key       string
}

type loginParams struct {
	Username string
	Password string
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenTtl    int64  `json:"tokenTtl"`
}
