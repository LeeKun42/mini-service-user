package request

type RegisterUser struct {
	Account  string `json:"account"`
	NickName string `json:"nick_name"`
	Passwd   string `json:"passwd"`
}

type UserLogin struct {
	Account string `url:"account"`
	Passwd  string `url:"passwd"`
}

type UserSearch struct {
	PageIndex int `url:"page_index"`
	PageSize  int `url:"page_size"`
}
