package echox

import (
	`net/http`
	`reflect`

	`github.com/labstack/echo/v4`
)

func (g *Group) RestfulPost(path string, handler restfulHandler, opts ...groupOption) *Route {
	_options := defaultGroupOptions()
	for _, opt := range opts {
		opt.applyGroup(_options)
	}

	return g.restful(
		http.MethodPost, path, handler,
		http.StatusCreated, http.StatusBadRequest,
		_options,
	)
}

func (g *Group) RestfulGet(path string, handler restfulHandler, opts ...groupOption) *Route {
	_options := defaultGroupOptions()
	for _, opt := range opts {
		opt.applyGroup(_options)
	}

	return g.restful(
		http.MethodGet, path, handler,
		http.StatusOK, http.StatusNoContent,
		_options,
	)
}

func (g *Group) RestfulPut(path string, handler restfulHandler, opts ...groupOption) *Route {
	_options := defaultGroupOptions()
	for _, opt := range opts {
		opt.applyGroup(_options)
	}

	return g.restful(
		http.MethodPut, path, handler,
		http.StatusOK, http.StatusBadRequest,
		_options,
	)
}

func (g *Group) RestfulDelete(path string, handler restfulHandler, opts ...groupOption) *Route {
	_options := defaultGroupOptions()
	for _, opt := range opts {
		opt.applyGroup(_options)
	}

	return g.restful(
		http.MethodDelete, path, handler,
		http.StatusNoContent, http.StatusBadRequest,
		_options,
	)
}

func (g *Group) restful(method string, path string, handler restfulHandler, successCode int, failedCode int, options *groupOptions) *Route {
	return &Route{
		Route: g.proxy.Add(method, path, func(ctx echo.Context) (err error) {
			var rsp interface{}
			if rsp, err = handler(parseContext(ctx)); nil != err {
				return
			}

			if g.checkFailed(rsp) {
				err = ctx.NoContent(failedCode)
			} else {
				options.code = successCode
				err = data(ctx, rsp, options.httpOptions)
			}

			return
		}, parseMiddlewares(options.middlewares...)...),
	}
}

func (g *Group) checkFailed(rsp interface{}) bool {
	return nil == rsp || reflect.ValueOf(rsp).IsZero()
}
