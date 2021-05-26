package echox

import (
	`time`
)

type options struct {
	// 地址
	addr string
	// 是否需要使用绑定功能
	binder bool
	// 是否需要使用数据验证功能
	validate bool
	// 是否显示徽标
	banner bool
	// 上下文路径，同一个机器上同一个端口运行不同的应用时需要加一个前级来区分不同的应用
	context string
	// 初始化方法
	inits []initFunc
	// 初始化路由方法
	routes []routeFunc
	// 错误处理
	error errorHandler
	// 退出超时时间
	shutdownTimeout time.Duration
	// 跨域
	crosEnable bool
	cros       crosConfig

	// 各种中间件
	// 打印堆栈信息
	panicStack panicStackConfig
	// Jwt配置
	jwtEnable bool
	jwt       JwtConfig
	// Http签名验证
	signatureEnable bool
	signature       signatureConfig
	// Casbin权限验证
	casbinEnable bool
	casbin       casbinConfig
}

func defaultOptions() *options {
	return &options{
		addr:            ":1323",
		binder:          true,
		validate:        true,
		banner:          false,
		error:           errorHandlerFunc,
		shutdownTimeout: 30 * time.Second,
		crosEnable:      true,
		cros: crosConfig{
			origins:     []string{"*"},
			credentials: true,
		},
		panicStack: panicStackConfig{
			size:              4 << 10,
			disableStackAll:   false,
			disablePrintStack: false,
		},
		jwtEnable:       false,
		signatureEnable: false,
		casbinEnable:    false,
	}
}
