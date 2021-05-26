package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

var methods = [...]string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	echo.PROPFIND,
	http.MethodPut,
	http.MethodTrace,
	echo.REPORT,
}
