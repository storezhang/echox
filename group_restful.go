package echox

import (
	`net/http`
	`reflect`

	`github.com/labstack/echo/v4`
)

func (g *Group) RestfulPost(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(
		http.MethodPost, path, handler,
		http.StatusCreated, http.StatusBadRequest,
		middlewares...,
	)
}

func (g *Group) RestfulGet(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(
		http.MethodGet, path, handler,
		http.StatusOK, http.StatusNoContent,
		middlewares...,
	)
}

func (g *Group) RestfulPut(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(
		http.MethodPut, path, handler,
		http.StatusOK, http.StatusBadRequest,
		middlewares...,
	)
}

func (g *Group) RestfulDelete(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(
		http.MethodDelete, path, handler,
		http.StatusNoContent, http.StatusBadRequest,
		middlewares...,
	)
}

func (g *Group) restful(method string, path string, handler restfulHandler, successCode int, failedCode int, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.proxy.Add(method, path, func(ctx echo.Context) (err error) {
			var rsp interface{}
			if rsp, err = handler(ctx.(*Context)); nil != err {
				return
			}

			if g.checkFailed(rsp) {
				err = ctx.NoContent(failedCode)
			} else {
				err = ctx.JSON(successCode, rsp)
			}

			return
		}, parseMiddlewares(middlewares...)...),
	}
}

func (g *Group) checkFailed(rsp interface{}) bool {
	return nil == rsp || reflect.ValueOf(rsp).IsZero()
}
