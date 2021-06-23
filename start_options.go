package echox

import (
	`time`
)

type startOptions struct {
	// 退出超时时间
	shutdownTimeout time.Duration
	// 初始化路由方法
	routes []routeFunc
}

func defaultStartOptions() *startOptions {
	return &startOptions{
		shutdownTimeout: 30 * time.Second,
		routes: []routeFunc{func(group *Group) {
			group.Get("/routes", routeHandler)
		}},
	}
}
