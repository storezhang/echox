package echox

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// PanicStackConfig Panic中间件配置
	PanicStackConfig struct {
		// Skipper 确定是不是要走中间件
		Skipper middleware.Skipper

		// StackSize 方法栈大小
		// 默认4KB
		StackSize int `yaml:"stack_size"`

		// DisableStackAll 是否禁止显示所有的栈信息
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack 禁止打印栈信息
		DisablePrintStack bool `yaml:"disable_print_stack"`
	}
)

var (
	// DefaultPanicStackConfig 默认配置
	DefaultPanicStackConfig = PanicStackConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
)

// PanicStack 中间件便捷方法
func PanicStack() echo.MiddlewareFunc {
	return PanicStackWithConfig(DefaultPanicStackConfig)
}

// PanicStackWithConfig 按配置生成中间件
func PanicStackWithConfig(config PanicStackConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultPanicStackConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultPanicStackConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						fmt.Printf("[Panic Stack Error] %v %s\n", err, stack[:length])
					}
					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
