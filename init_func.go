package echox

import (
	`github.com/labstack/echo/v4`
)

type initFunc func(echo *echo.Echo)
