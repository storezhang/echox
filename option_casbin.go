package echox

import (
	`github.com/casbin/casbin/v2`
	`github.com/labstack/echo/v4/middleware`
)

var _ option = (*optionCasbin)(nil)

type optionCasbin struct {
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

// Casbin Http签名
func Casbin(enforcer *casbin.Enforcer, jwt JwtConfig, source roleSource) *optionCasbin {
	return CasbinWithConfig(middleware.DefaultSkipper, enforcer, jwt, false, source)
}

// CasbinWithConfig Http签名
func CasbinWithConfig(skipper middleware.Skipper, enforcer *casbin.Enforcer, jwt JwtConfig, trailingSlash bool, source roleSource) *optionCasbin {
	return &optionCasbin{
		skipper:       skipper,
		enforcer:      enforcer,
		jwt:           jwt,
		trailingSlash: trailingSlash,
		source:        source,
	}
}

func (j *optionCasbin) apply(options *options) {
	options.casbin.skipper = j.skipper
	options.casbin.enforcer = j.enforcer
	options.casbin.jwt = j.jwt
	options.casbin.trailingSlash = j.trailingSlash
	options.casbin.source = j.source
	options.casbinEnable = true
}
