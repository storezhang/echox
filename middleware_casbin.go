package echox

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// roleSource 获得用户的角色编号列表
type roleSource interface {
	// GetsRoleId 获得用户的角色编号列表
	GetsRoleId(user interface{}) (ids []int64, err error)
}

func casbinFunc(config casbinConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if config.skipper(ctx) {
				err = next(ctx)

				return
			}

			var pass bool
			if pass, err = config.checkPermission(ctx); err == nil && pass {
				err = next(ctx)
			} else if err != nil {
				err = echo.NewHTTPError(http.StatusForbidden, err.Error())
			} else {
				err = echo.ErrForbidden
			}

			return
		}
	}
}
