package echox

var _ option = (*optionRoutes)(nil)

type optionRoutes struct {
	// 路由器
	routes []routeFunc
}

// Routes 配置路由器
func Routes(routes ...routeFunc) *optionRoutes {
	return &optionRoutes{routes: routes}
}

func (r *optionRoutes) apply(options *options) {
	options.routes = append(options.routes, r.routes...)
}
