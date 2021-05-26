package echox

import (
	`fmt`
	`runtime`

	`github.com/labstack/echo/v4`
)

func panicStackFunc(config panicStackConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.size)
					length := runtime.Stack(stack, !config.disableStackAll)
					if !config.disablePrintStack {
						fmt.Printf("[异常堆栈信息] %v %s\n", err, stack[:length])
					}
					ctx.Error(err)
				}
			}()

			return next(ctx)
		}
	}
}
