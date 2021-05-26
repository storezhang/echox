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

type casbinConfig struct {
	// 确定是不是要走中间件
	skipper middleware.Skipper
	// Casbin的权限验证模块
	enforcer *casbin.Enforcer
	// Jwt的配置
	jwt JwtConfig
	// 是否包含尾部斜杠
	trailingSlash bool
	// 用户角色权限
	source roleSource
}

func (cc *casbinConfig) checkPermission(ctx echo.Context) (pass bool, err error) {
	var (
		user    gox.BaseUser
		roleIds []int64
		ec      = Context{
			Context: ctx,
			jwt:     cc.jwt,
		}
	)

	if err = ec.Subject(user); nil != err {
		return
	}

	if roleIds, err = cc.source.GetsRoleId(user.Id); nil != err {
		return
	}

	path := ctx.Request().URL.Path
	if pass, err = cc.checkCasbinPermission(path, methodMapping[ctx.Request().Method], roleIds...); nil != err {
		return
	}

	// 取得Path
	// 统一加上最后的斜杠
	if !pass && cc.trailingSlash {
		path += "/"
		pass, err = cc.checkCasbinPermission(path, methodMapping[ctx.Request().Method], roleIds...)
	}

	return
}

func (cc *casbinConfig) checkCasbinPermission(obj string, act string, roleIds ...int64) (pass bool, err error) {
	for _, roleId := range roleIds {
		roleIdStr := strconv.FormatInt(roleId, 10)
		// 调用Casbin检查权限
		if pass, err = cc.enforcer.Enforce(roleIdStr, obj, act); nil != err {
			break
		}

		// 已经有权限，提前结束
		if pass {
			break
		}
	}

	return
}
