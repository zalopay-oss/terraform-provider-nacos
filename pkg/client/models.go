package client

type Configuration struct {
	Namespace   string `http:"tenant,query,form" json:"tenant"`
	Group       string `http:"group,query,form" json:"group"`
	Key         string `http:"dataId,query,form" json:"dataId"`
	Value       string `http:"content,query,form" json:"content"`
	Description string `http:"desc,query,form" json:"desc"`
}

type ConfigurationId struct {
	Namespace string `http:"tenant,query,form"`
	Group     string `http:"group,query,form"`
	Key       string `http:"dataId,query,form"`
}

type authParams struct {
	AccessToken string `http:"accessToken,query"`
}

type optionalParams struct {
	Show string `http:"show,query"`
}

type loginParams struct {
	Username string `http:"username,form"`
	Password string `http:"password,form"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenTtl    int64  `json:"tokenTtl"`
}
