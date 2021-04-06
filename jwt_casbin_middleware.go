package echox

import (
	`encoding/json`
	"net/http"
	"strconv"

	"github.com/storezhang/gox"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// JWTCasbinConfig Casbin配置
	JWTCasbinConfig struct {
		// Skipper 确定是不是要走中间件
		Skipper middleware.Skipper
		// Enforcer Casbin的权限验证模块
		Enforcer *casbin.Enforcer
		// JWT JWT的配置
		JWT *JWTConfig
		// TrailingSlash 是否包含尾部斜杠
		TrailingSlash bool
		// RoleSource 用户角色权限
		RoleSource RoleSource
	}

	// RoleSource 获得用户的角色编号列表
	RoleSource interface {
		// GetsRoleIds 获得用户的角色编号列表
		GetsRoleId(user interface{}) ([]int64, error)
	}
)

var (
	DefaultJWTCasbinConfig = JWTCasbinConfig{
		Skipper: middleware.DefaultSkipper,
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

func (jcc JWTCasbinConfig) String() string {
	jsonBytes, _ := json.MarshalIndent(jcc, "", "    ")

	return string(jsonBytes)
}

func (jcc *JWTCasbinConfig) CheckPermission(c echo.Context) (checked bool, err error) {
	var (
		user    gox.BaseUser
		roleIds []int64
		ec      = EchoContext{
			Context: c,
			jwt:     jcc.JWT,
		}
	)

	if err = ec.Subject(user); nil != err {
		return
	}

	if roleIds, err = jcc.RoleSource.GetsRoleId(user.Id); nil != err {
		return
	}

	path := c.Request().URL.Path
	if checked, err = jcc.checkPermission(path, methodMapping[c.Request().Method], roleIds...); nil != err {
		return
	}

	// 取得Path
	// 统一加上最后的斜杠
	if !checked && jcc.TrailingSlash {
		path += "/"
		checked, err = jcc.checkPermission(path, methodMapping[c.Request().Method], roleIds...)
	}

	return
}

func (jcc *JWTCasbinConfig) checkPermission(ojb string, act string, roleIds ...int64) (checked bool, err error) {
	for _, roleId := range roleIds {
		roleIdStr := strconv.FormatInt(roleId, 10)
		// 调用Casbin检查权限
		if checked, err = jcc.Enforcer.Enforce(roleIdStr, ojb, act); nil != err {
			break
		}

		// 已经有权限，提前结束
		if checked {
			break
		}
	}

	return
}
