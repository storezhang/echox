package echox

import (
	`net/http`

	`github.com/labstack/echo/v4`
)

type (
	httpOption interface {
		applyHttp(options *httpOptions)
	}

	httpOptions struct {
		code        int
		contentType string
		serializer  serializer
	}
)

func defaultHttpOptions() *httpOptions {
	return &httpOptions{
		code:        http.StatusOK,
		contentType: echo.MIMEApplicationJSONCharsetUTF8,
		serializer:  serializerJson,
	}
}
