package controller

type Response struct {
	//失败是0，成功是1
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
