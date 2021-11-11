package echox

import (
	`context`
	`fmt`
	`os`
	`os/signal`
	`strings`

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
	_options := defaultOptions
	for _, opt := range opts {
		opt.apply(_options)
	}

	// 创建Echo服务器
	server := echo.New()
	server.HideBanner = !_options.banner

	// 初始化
	for _, init := range _options.inits {
		init(server)
	}

	// 数据验证
	if _options.validate {
		server.Validator = &validate{validate: validatorx.New()}
	}

	// 初始化绑定
	if nil != _options.binder {
		server.Binder = _options.binder
	}

	// 处理错误
	server.HTTPErrorHandler = echo.HTTPErrorHandler(_options.error)

	// 初始化中间件
	server.Pre(middleware.MethodOverride())
	server.Pre(middleware.RemoveTrailingSlash())

	// server.Use(middleware.CSRF())
	server.Use(middleware.Logger())
	server.Use(middleware.RequestID())
	// 配置跨域
	if _options.crosEnable {
		cors := middleware.DefaultCORSConfig
		cors.AllowMethods = append(cors.AllowMethods, string(gox.HttpMethodOptions))
		cors.AllowOrigins = _options.cros.origins
		cors.AllowCredentials = _options.cros.credentials
		server.Use(middleware.CORSWithConfig(cors))
	}

	// 打印堆栈信息
	// 方便调试，默认处理没有换行，很难内眼查看堆栈信息
	server.Use(panicStackFunc(_options.panicStack))

	// 增加自定义上下文
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			return next(&Context{
				Context: ctx,
			})
		}
	})

	return &Echo{
		Echo:    server,
		options: _options,
	}
}

func (e *Echo) Start(opts ...startOption) (err error) {
	_options := defaultStartOptions()
	for _, opt := range opts {
		opt.applyStart(_options)
	}

	// 处理路由
	if 0 != len(_options.routes) {
		group := &Group{proxy: e.Group(e.options.context)}
		for _, route := range _options.routes {
			route(group)
		}
	}

	// 在另外的协程中启动服务器，实现优雅地关闭（Graceful Shutdown）
	if _options.graceful {
		go func() {
			err = e.graceful(_options)
		}()
	} else {
		err = e.Echo.Start(e.options.addr)
	}

	return
}

func (e *Echo) Shutdown(opts ...stopOption) error {
	_options := defaultStopOptions()
	for _, opt := range opts {
		opt.applyStop(_options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), _options.timeout)
	defer cancel()

	return e.Echo.Shutdown(ctx)
}

func (e *Echo) Url(api string) (url string) {
	if strings.HasPrefix(api, `/`) {
		api = api[1:]
	}

	if `` != e.options.proxy {
		url = fmt.Sprintf(`%s/%s`, e.options.proxy, api)
	} else {
		url = fmt.Sprintf(`%s/%s`, e.options.addr, api)
	}

	return
}

func (e *Echo) graceful(options *startOptions) (err error) {
	if err = e.Echo.Start(e.options.addr); nil != err {
		return
	}

	// 等待系统退出中断并响应
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), options.shutdownTimeout)
	defer cancel()
	err = e.Echo.Shutdown(ctx)

	return
}
