package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

// Group 分组，模拟echo.Group，并增加Context转换
type Group struct {
	group *echo.Group
}

func (g *Group) Use(middlewares ...MiddlewareFunc) {
	g.group.Use(parseMiddlewares(middlewares...)...)
}

func (g *Group) Connect(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodConnect, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Delete(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodDelete, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Get(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodGet, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Head(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodHead, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Options(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodOptions, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Patch(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPatch, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Post(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPost, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Put(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodPut, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Trace(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodTrace, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Any(path string, handler handlerFunc, middlewares ...MiddlewareFunc) (routes []*Route) {
	routes = make([]*Route, len(methods))
	for index, method := range methods {
		routes[index] = g.Add(method, path, handler, middlewares...)
	}

	return
}

func (g *Group) Match(methods []string, path string, handler handlerFunc, middlewares ...MiddlewareFunc) (routes []*Route) {
	routes = make([]*Route, len(methods))
	for index, method := range methods {
		routes[index] = g.Add(method, path, handler, middlewares...)
	}

	return
}

func (g *Group) Group(prefix string, middlewares ...MiddlewareFunc) (ag *Group) {
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

func (g *Group) Add(method string, path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(method, path, func(ctx echo.Context) error {
			return handler(ctx.(*Context))
		}, parseMiddlewares(middlewares...)...),
	}
}
