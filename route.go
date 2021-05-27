package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

type (
	routeFunc func(group *Group)

	// Route 描述一个路由器，快捷方式，方便用户操作
	Route struct {
		*echo.Route
	}
)

func routeHandler(ctx *Context) error {
	return ctx.JSON(http.StatusOK, ctx.Echo().Routes())
}
