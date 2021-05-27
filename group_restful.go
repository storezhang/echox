package echox

import (
	`net/http`
	`reflect`

	`github.com/labstack/echo/v4`
)

func (g *Group) RestfulPost(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(path, handler, http.StatusCreated, http.StatusBadRequest, middlewares...)
}

func (g *Group) RestfulGet(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(path, handler, http.StatusOK, http.StatusNotFound, middlewares...)
}

func (g *Group) RestfulUpdate(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(path, handler, http.StatusOK, http.StatusBadRequest, middlewares...)
}

func (g *Group) RestfulDelete(path string, handler restfulHandler, middlewares ...MiddlewareFunc) *Route {
	return g.restful(path, handler, http.StatusNoContent, http.StatusBadRequest, middlewares...)
}

func (g *Group) restful(path string, handler restfulHandler, successCode int, failedCode int, middlewares ...MiddlewareFunc) *Route {
	return &Route{
		Route: g.group.Add(http.MethodGet, path, func(ctx echo.Context) (err error) {
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
