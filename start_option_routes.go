package echox

var _ startOption = (*startOptionRoutes)(nil)

type startOptionRoutes struct {
	// 路由器
	routes []routeFunc
}

// Routes 配置路由器
func Routes(routes ...routeFunc) *startOptionRoutes {
	return &startOptionRoutes{routes: routes}
}

func (r *startOptionRoutes) applyStart(startOptions *startOptions) {
	startOptions.routes = append(startOptions.routes, r.routes...)
}
