package echox

var _ groupOption = (*optionMiddlewares)(nil)

type optionMiddlewares struct {
	middlewares []MiddlewareFunc
}

// Middleware 配置中间件
func Middleware(middleware MiddlewareFunc) *optionMiddlewares {
	return &optionMiddlewares{
		middlewares: []MiddlewareFunc{middleware},
	}
}

// Middlewares 配置中间件
func Middlewares(middlewares ...MiddlewareFunc) *optionMiddlewares {
	return &optionMiddlewares{
		middlewares: middlewares,
	}
}

func (m *optionMiddlewares) applyGroup(options *groupOptions) {
	options.middlewares = append(options.middlewares, m.middlewares...)
}
