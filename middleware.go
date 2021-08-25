package echox

import (
	`github.com/labstack/echo/v4`
)

// MiddlewareFunc 中间件
type MiddlewareFunc func(handlerFunc) handlerFunc

func parseMiddlewares(middlewares ...MiddlewareFunc) (ems []echo.MiddlewareFunc) {
	length := len(middlewares)
	if 0 == length {
		return
	}

	ems = make([]echo.MiddlewareFunc, 0, length)
	for _, middleware := range middlewares {
		ems = append(ems, func(next echo.HandlerFunc) echo.HandlerFunc {
			handler := middleware(func(ctx *Context) (err error) {
				return next(ctx.Context)
			})

			return func(ctx echo.Context) error {
				return handler(ctx.(*Context))
			}
		})
	}

	return
}
