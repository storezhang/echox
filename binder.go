package echox

import (
	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type binder struct{}

func (b *binder) Bind(data interface{}, ctx echo.Context) (err error) {
	// 处理默认值
	defaults.SetDefaults(data)

	db := new(echo.DefaultBinder)
	if err = db.Bind(data, ctx); err != echo.ErrUnsupportedMediaType {
		return
	}

	return
}
