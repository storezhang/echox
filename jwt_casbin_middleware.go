package echox

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/storezhang/gox"

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
		// URL前缀
		UrlPrefix string
		// 用户角色权限
		RoleSource RoleSource
	}

	RoleSource interface {
		GetsUserRoleIds(int64) ([]int64, error)
	}
)

var (
	DefaultJWTCasbinConfig = JWTCasbinConfig{
		Skipper:   middleware.DefaultSkipper,
		UrlPrefix: "/api",
	}

	methodMapping = map[string]string{
		"GET":    "r",
		"POST":   "c",
		"PUT":    "u",
		"DELETE": "d",
		"*":      "*",
	}
)

func JWTCasbinMiddleware(jwtCasbinCfg JWTCasbinConfig, e *casbin.Enforcer, jwt *JWTConfig, trailingSlash bool) echo.MiddlewareFunc {
	c := jwtCasbinCfg
	c.Enforcer = e
	c.JWT = jwt
	c.TrailingSlash = trailingSlash

	return JWTCasbinWithConfig(c)
}

func JWTCasbinWithConfig(config JWTCasbinConfig) echo.MiddlewareFunc {
	if nil == config.Skipper {
		config.Skipper = DefaultJWTCasbinConfig.Skipper
	}
	if "" == strings.TrimSpace(config.UrlPrefix) {
		config.UrlPrefix = DefaultJWTCasbinConfig.UrlPrefix
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

func (jcc *JWTCasbinConfig) CheckPermission(c echo.Context) (checked bool, err error) {
	var (
		user    gox.BaseUser
		roleIds []int64
		ec      = EchoContext{
			Context: c,
			JWT:     jcc.JWT,
		}
	)

	path := c.Request().URL.Path
	// 取得Path
	// 统一加上最后的斜杠
	if jcc.TrailingSlash && 0 == strings.LastIndex(path[len(jcc.UrlPrefix):], "/") {
		path += "/"
	}

	if user, err = ec.User(); nil != err {
		return
	}

	if roleIds, err = jcc.RoleSource.GetsUserRoleIds(user.Id); nil != err {
		return
	}
	for _, roleId := range roleIds {
		roleIdStr := strconv.FormatInt(roleId, 10)
		// 调用Casbin检查权限
		if checked, err = jcc.Enforcer.Enforce(roleIdStr, path, methodMapping[c.Request().Method]); nil != err {
			return
		} else if checked {
			return
		}
	}

	return
}
