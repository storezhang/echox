package echox

import (
	`strconv`

	`github.com/casbin/casbin/v2`
	`github.com/labstack/echo/v4`
	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
)

var methodMapping = map[string]string{
	"GET":    "r",
	"POST":   "c",
	"PUT":    "u",
	"DELETE": "d",
	"*":      "*",
}

type Casbin struct {
	// 确定是不是要走中间件
	skipper middleware.Skipper
	// Casbin的权限验证模块
	enforcer *casbin.Enforcer
	// Jwt的配置
	jwt Jwt
	// 是否包含尾部斜杠
	trailingSlash bool
	// 用户角色权限
	source roleSource
}

// NewCasbin Casbin权限验证
func NewCasbin(enforcer *casbin.Enforcer, jwt Jwt, source roleSource) *Casbin {
	return NewCasbinWithConfig(middleware.DefaultSkipper, enforcer, jwt, false, source)
}

// NewCasbinWithConfig Casbin权限验证
func NewCasbinWithConfig(
	skipper middleware.Skipper,
	enforcer *casbin.Enforcer,
	jwt Jwt,
	trailingSlash bool,
	source roleSource,
) *Casbin {
	return &Casbin{
		skipper:       skipper,
		enforcer:      enforcer,
		jwt:           jwt,
		trailingSlash: trailingSlash,
		source:        source,
	}
}

func (c *Casbin) checkPermission(ctx echo.Context) (pass bool, err error) {
	var (
		user    gox.BaseUser
		roleIds []int64
	)

	if err = c.jwt.Subject(&Context{Context: ctx}, user); nil != err {
		return
	}

	if roleIds, err = c.source.GetsRoleId(user.Id); nil != err {
		return
	}

	path := ctx.Request().URL.Path
	if pass, err = c.checkCasbinPermission(path, methodMapping[ctx.Request().Method], roleIds...); nil != err {
		return
	}

	// 取得Path
	// 统一加上最后的斜杠
	if !pass && c.trailingSlash {
		path += "/"
		pass, err = c.checkCasbinPermission(path, methodMapping[ctx.Request().Method], roleIds...)
	}

	return
}

func (c *Casbin) checkCasbinPermission(obj string, act string, roleIds ...int64) (pass bool, err error) {
	for _, roleId := range roleIds {
		roleIdStr := strconv.FormatInt(roleId, 10)
		// 调用Casbin检查权限
		if pass, err = c.enforcer.Enforce(roleIdStr, obj, act); nil != err {
			break
		}

		// 已经有权限，提前结束
		if pass {
			break
		}
	}

	return
}
