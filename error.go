package echox

// Error 接口，符合条件的错误统一处理
type Error interface {
	// ErrorCode 返回错误码
	ErrorCode() int
	// Message 返回错误消息
	Message() string
	// Data 返回错误实体
	// 在某些错误下，可能需要返回额外的信息给前端处理
	// 比如，认证错误，需要返回哪些字段有错误
	Data() interface{}
}
