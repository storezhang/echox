package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

// Group 分组，模拟echo.Group，并增加Context转换
type Group struct {
	group *echo.Group
}

func (g *Group) Use(middlewares ...middlewareFunc) {
	g.group.Use(parseMiddlewares(middlewares...)...)
}

func (g *Group) CONNECT(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodConnect, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) DELETE(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodDelete, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) GET(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodGet, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) HEAD(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodHead, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) OPTIONS(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodOptions, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) PATCH(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPatch, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) POST(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPost, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) PUT(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPut, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) TRACE(path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodTrace, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Any(path string, handler handlerFunc, middlewares ...middlewareFunc) (routes []*Route) {
	routes = make([]*Route, len(methods))
	for index, method := range methods {
		routes[index] = g.Add(method, path, handler, middlewares...)
	}

	return
}

func (g *Group) Match(methods []string, path string, handler handlerFunc, middlewares ...middlewareFunc) (routes []*Route) {
	routes = make([]*Route, len(methods))
	for index, method := range methods {
		routes[index] = g.Add(method, path, handler, middlewares...)
	}

	return
}

func (g *Group) Group(prefix string, middlewares ...middlewareFunc) (ag *Group) {
	return &Group{
		group: g.group.Group(prefix, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Static(prefix string, root string) {
	g.group.Static(prefix, root)
}

func (g *Group) File(path string, file string) {
	g.group.File(path, file)
}

func (g *Group) Add(method string, path string, handler handlerFunc, middlewares ...middlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(method, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}
