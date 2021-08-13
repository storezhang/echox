package echox

import (
	`time`
)

type startOptions struct {
	*options

	graceful        bool
	shutdownTimeout time.Duration
	routes          []routeFunc
}

func defaultStartOptions() *startOptions {
	return &startOptions{
		options: defaultOptions,

		graceful:        false,
		shutdownTimeout: 30 * time.Second,
		routes: []routeFunc{func(group *Group) {
			group.Get("/routes", routeHandler).Name = "所有路由信息"
		}},
	}
}
