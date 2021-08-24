package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

// Group 分组，模拟echo.Group，并增加Context转换
// 使用代理设计模式
type Group struct {
	proxy *echo.Group
}

func (g *Group) Use(middlewares ...MiddlewareFunc) {
	g.proxy.Use(parseMiddlewares(middlewares...)...)
}

func (g *Group) Connect(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodConnect, path, handler, middlewares...)
}

func (g *Group) Delete(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodDelete, path, handler, middlewares...)
}

func (g *Group) Get(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodGet, path, handler, middlewares...)
}

func (g *Group) Head(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodHead, path, handler, middlewares...)
}

func (g *Group) Options(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodOptions, path, handler, middlewares...)
}

func (g *Group) Patch(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPatch, path, handler, middlewares...)
}

func (g *Group) Post(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPost, path, handler, middlewares...)
}

func (g *Group) Put(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPut, path, handler, middlewares...)
}

func (g *Group) Trace(path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return g.Add(http.MethodTrace, path, handler, middlewares...)
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
		proxy: g.proxy.Group(prefix, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) Static(prefix string, root string) {
	g.proxy.Static(prefix, root)
}

func (g *Group) File(path string, file string) {
	g.proxy.File(path, file)
}

func (g *Group) Add(method string, path string, handler handlerFunc, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.proxy.Add(method, path, func(ctx echo.Context) (err error) {
			if _ctx, ok := ctx.(*Context); ok {
				err = handler(_ctx)
			} else {
				err = handler(&Context{Context: ctx})
			}

			return
		}, parseMiddlewares(middlewares...)...),
	}
}
