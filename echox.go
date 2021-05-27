package echox

import (
	`context`
	`os`
	`os/signal`

	`github.com/labstack/echo/v4`
	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
	`github.com/storezhang/validatorx`
)

type initFunc func(echo *echo.Echo)

// Start 启动服务
func Start(opts ...option) (err error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	// 创建Echo服务器
	server := echo.New()
	server.HideBanner = !options.banner

	// 初始化
	for _, init := range options.inits {
		init(server)
	}

	// 处理路由
	if 0 != len(options.routes) {
		group := &Group{proxy: server.Group(options.context)}
		for _, route := range options.routes {
			route(group)
		}
	}

	// 数据验证
	if options.validate {
		server.Validator = &validate{validate: validatorx.New()}
	}

	// 初始化绑定
	if options.binder {
		server.Binder = &binder{}
	}

	// 处理错误
	server.HTTPErrorHandler = echo.HTTPErrorHandler(options.error)

	// 初始化中间件
	server.Pre(middleware.MethodOverride())
	server.Pre(middleware.RemoveTrailingSlash())

	// server.Use(middleware.CSRF())
	server.Use(middleware.Logger())
	server.Use(middleware.RequestID())
	// 配置跨域
	if options.crosEnable {
		cors := middleware.DefaultCORSConfig
		cors.AllowMethods = append(cors.AllowMethods, string(gox.HttpMethodOptions))
		cors.AllowOrigins = options.cros.origins
		cors.AllowCredentials = options.cros.credentials
		server.Use(middleware.CORSWithConfig(cors))
	}

	// 增加内置中间件
	// 增加Jwt
	if options.jwtEnable {
		server.Use(jwtFunc(options.jwt))
	}
	// 增加Http签名验证
	if options.signatureEnable {
		server.Use(signatureFunc(options.signature))
	}
	// 增加权限验证
	if options.casbinEnable {
		server.Use(casbinFunc(options.casbin))
	}
	// 打印堆栈信息
	// 方便调试，默认处理没有换行，很难内眼查看堆栈信息
	server.Use(panicStackFunc(options.panicStack))

	// 增加自定义上下文
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Context{
				Context: c,
				jwt:     options.jwt,
			})
		}
	})

	// 在另外的协程中启动服务器，实现优雅地关闭（Graceful Shutdown）
	go func() {
		err = server.Start(options.addr)
	}()

	// 等待系统退出中断并响应
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), options.shutdownTimeout)
	defer cancel()
	err = server.Shutdown(ctx)

	return
}
