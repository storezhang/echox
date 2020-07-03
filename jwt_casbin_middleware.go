package echox

import (
	"net/http"
	`strings`

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	JWTCasbinConfig struct {
		// 确定是不是要走中间件
		Skipper middleware.Skipper
		// Casbin的权限验证模块
		Enforcer *casbin.Enforcer
		// JWT的配置
		JWT *JWTConfig
		// 是否包含尾部斜杠
		TrailingSlash bool
	}
)

var (
	DefaultJWTCasbinConfig = JWTCasbinConfig{
		Skipper: middleware.DefaultSkipper,
	}

	MethodMapping = map[string]string{
		"GET":    "r",
		"POST":   "c",
		"PUT":    "u",
		"DELETE": "d",
		"*":      "*",
	}
)

func JWTCasbinMiddleware(e *casbin.Enforcer, jwt *JWTConfig, trailingSlash bool) echo.MiddlewareFunc {
	c := DefaultJWTCasbinConfig
	c.Enforcer = e
	c.JWT = jwt
	c.TrailingSlash = trailingSlash

	return JWTCasbinWithConfig(c)
}

func JWTCasbinWithConfig(config JWTCasbinConfig) echo.MiddlewareFunc {
	if nil == config.Skipper {
		config.Skipper = DefaultJWTCasbinConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if pass, err := config.CheckPermission(c); err == nil && pass {
				return next(c)
			} else if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}

			return echo.ErrForbidden
		}
	}
}

func (jcc *JWTCasbinConfig) CheckPermission(c echo.Context) (bool, error) {
	ec := EchoContext{
		Context: c,
		JWT:     jcc.JWT,
	}

	path := c.Request().URL.Path
	// 取得Path
	// 统一加上最后的斜杠
	if jcc.TrailingSlash && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if user, err := ec.User(); nil != err {
		return false, err
	} else {
		return jcc.Enforcer.Enforce(user.IdString(), path, MethodMapping[c.Request().Method])
	}
}
