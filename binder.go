package echox

import (
	`reflect`

	`github.com/labstack/echo/v4`
	`github.com/mcuadros/go-defaults`
)

type binder struct{}

func (b *binder) Bind(data interface{}, ctx echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(data, ctx); err != echo.ErrUnsupportedMediaType {
		return
	}

	// 处理默认值
	// 区分指针类型和非指针类型
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		defaults.SetDefaults(data)
	} else {
		defaults.SetDefaults(&data)
	}

	return
}
