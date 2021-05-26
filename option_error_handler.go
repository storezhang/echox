package echox

var _ option = (*optionErrorHandler)(nil)

type optionErrorHandler struct {
	// 错误处理
	handler errorHandler
}

// ErrorHandler 绑定地址
func ErrorHandler(handler errorHandler) *optionErrorHandler {
	return &optionErrorHandler{
		handler: handler,
	}
}

func (e *optionErrorHandler) apply(options *options) {
	options.error = e.handler
}
