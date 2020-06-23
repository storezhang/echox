package echox

import (
	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type DefaultValueBinder struct{}

func (dvb *DefaultValueBinder) Bind(i interface{}, c echo.Context) (err error) {
	defaults.SetDefaults(i)

	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	return
}
