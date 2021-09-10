package echox

var _ groupOption = (*optionMiddleware)(nil)

type optionMiddleware struct {
	middleware MiddlewareFunc
}

// Middleware 配置中间件
func Middleware(middleware MiddlewareFunc) *optionMiddleware {
	return &optionMiddleware{
		middleware: middleware,
	}
}

func (m *optionMiddleware) applyGroup(options *groupOptions) {
	options.middlewares = append(options.middlewares, m.middleware)
}
