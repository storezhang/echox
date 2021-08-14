package echox

import (
	`reflect`

	`github.com/labstack/echo/v4`
	`github.com/mcuadros/go-defaults`
)

type binder struct{}

func (b *binder) Bind(req interface{}, ctx echo.Context) (err error) {
	defaultBinder := new(echo.DefaultBinder)
	if err = defaultBinder.Bind(req, ctx); nil != err {
		return
	}

	// 处理默认值
	// 区分指针类型和非指针类型
	if reflect.ValueOf(req).Kind() == reflect.Ptr {
		defaults.SetDefaults(req)
	} else {
		defaults.SetDefaults(&req)
	}

	return
}
