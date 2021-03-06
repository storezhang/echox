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

// Echo 组织echo.Echo启动
type Echo struct {
	*echo.Echo

	options *options
}

func New(opts ...option) *Echo {
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

	// 打印堆栈信息
	// 方便调试，默认处理没有换行，很难内眼查看堆栈信息
	server.Use(panicStackFunc(options.panicStack))

	// 增加自定义上下文
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Context{
				Context: c,
			})
		}
	})

	return &Echo{
		Echo:    server,
		options: options,
	}
}

func (e *Echo) Start(opts ...startOption) (err error) {
	options := defaultStartOptions()
	for _, opt := range opts {
		opt.applyStart(options)
	}

	// 处理路由
	if 0 != len(options.routes) {
		group := &Group{proxy: e.Group(e.options.context)}
		for _, route := range options.routes {
			route(group)
		}
	}

	// 在另外的协程中启动服务器，实现优雅地关闭（Graceful Shutdown）
	go func() {
		err = e.Echo.Start(e.options.addr)
	}()

	// 等待系统退出中断并响应
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), options.shutdownTimeout)
	defer cancel()
	err = e.Shutdown(ctx)

	return
}
