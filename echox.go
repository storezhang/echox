package echox

import (
	`context`
	`fmt`
	`net/http`
	`os`
	`os/signal`
	`strconv`
	`strings`
	`time`

	`github.com/go-playground/validator/v10`
	`github.com/labstack/echo/v4`
	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
	`github.com/storezhang/validatorx`
)

var (
	DefaultEchoConfig = &EchoConfig{
		Ip:                 "",
		Port:               1323,
		BasePath:           "",
		Validate:           true,
		DefaultValueBinder: true,
		ErrorHandler:       true,
		JWT:                nil,
		Init:               nil,
		Routes:             nil,
	}
)

type (
	EchoFunc   func(e *echo.Echo)
	RouteFunc  func(g *echo.Group)
	EchoConfig struct {
		Ip                 string
		Port               int
		BasePath           string
		Validate           bool
		DefaultValueBinder bool
		ErrorHandler       bool
		JWT                *JWTConfig
		Init               EchoFunc
		Routes             []RouteFunc
	}
)

func (ec *EchoConfig) Address() (address string) {
	if "" != strings.TrimSpace(ec.Ip) {
		address = fmt.Sprintf("%s:%d", ec.Ip, ec.Port)
	} else {
		address = fmt.Sprintf(":%d", ec.Port)
	}

	return
}

func Start() {
	StartWith(DefaultEchoConfig)
}

func StartWith(ec *EchoConfig) {
	// 创建Echo对象
	e := echo.New()

	if nil != ec.Init {
		ec.Init(e)
	}
	if nil != ec.Routes {
		g := e.Group(ec.BasePath)
		for _, route := range ec.Routes {
			route(g)
		}
	}

	// 初始化Validator
	if ec.Validate {
		// 数据验证
		e.Validator = validatorx.New()
	}

	// 初始化绑定
	if ec.DefaultValueBinder {
		e.Binder = &DefaultValueBinder{}
	}

	// 处理错误
	if ec.ErrorHandler {
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			type response struct {
				ErrorCode int         `json:"errorCode"`
				Message   string      `json:"message"`
				Data      interface{} `json:"data"`
			}
			rsp := response{}

			statusCode := http.StatusInternalServerError
			switch re := err.(type) {
			case *echo.HTTPError:
				statusCode = re.Code
				rsp.ErrorCode = 9902
				rsp.Message = "处理请求失败"
				if nil != re.Internal {
					rsp.Data = re.Internal.Error()
				}
			case validator.ValidationErrors:
				statusCode = http.StatusBadRequest
				lang := c.Request().Header.Get(gox.HeaderAcceptLanguage)
				rsp.ErrorCode = 9901
				rsp.Message = "数据验证错误"
				rsp.Data = validatorx.I18n(lang, re)
			case *gox.CodeError:
				rsp.ErrorCode = int(re.ToErrorCode())
				rsp.Message = re.ToMessage()
				rsp.Data = re.ToData()
			default:
				rsp.ErrorCode = 9903
				rsp.Message = "服务器内部错误"
			}

			_ = c.JSON(statusCode, rsp)
		}
	}

	// 初始化中间件
	e.Pre(middleware.MethodOverride())
	e.Pre(middleware.RemoveTrailingSlash())

	// e.Use(middleware.CSRF())
	e.Use(middleware.Logger())
	// 打印堆栈信息
	// 方便调试，默认处理没有换行，很难内眼查看堆栈信息
	e.Use(PanicStack())
	e.Use(middleware.RequestID())

	// 符合JWT和Casbin的上下文
	if nil != ec.JWT {
		e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return h(NewContext(c, ec.JWT))
			}
		})
	}

	// 启动Server
	go func() {
		if err := e.Start(ec.Address()); nil != err {
			e.Logger.Fatal(err)
		}
	}()

	// 等待系统退出中断并响应
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); nil != err {
		e.Logger.Fatal(err)
	}
}

func Int64Param(c echo.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}

func IntParam(c echo.Context, name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}
