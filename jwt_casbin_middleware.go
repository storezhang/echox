package echox

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	JWTCasbinConfig struct {
		Skipper  middleware.Skipper
		Enforcer *casbin.Enforcer
		JWT      *JWTConfig
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

func JWTCasbinMiddleware(e *casbin.Enforcer, jwt *JWTConfig) echo.MiddlewareFunc {
	c := DefaultJWTCasbinConfig
	c.Enforcer = e
	c.JWT = jwt

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
	if user, err := ec.User(); nil != err {
		return false, err
	} else {
		return jcc.Enforcer.Enforce(user.IdString(), c.Request().URL.Path, MethodMapping[c.Request().Method])
	}
}
