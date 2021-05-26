package echox

type response struct {
	// 错误码
	ErrorCode int `json:"errorCode"`
	// 消息
	Message string `json:"message"`
	// 附加数据
	Data interface{} `json:"data"`
}
