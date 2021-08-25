package echox

import (
	`bytes`
	`net/http`
	`os`

	`github.com/labstack/echo/v4`
	`github.com/storezhang/gox`
)

const defaultIndent = "  "

// Context 自定义的Echo上下文
type Context struct {
	echo.Context
}

func (c *Context) Fill(data interface{}) (err error) {
	if err = c.Bind(data); nil != err {
		return
	}
	err = c.Validate(data)

	return
}

func (c *Context) HttpFile(file http.File) (err error) {
	defer func() {
		_ = file.Close()
	}()

	var info os.FileInfo
	if info, err = file.Stat(); nil != err {
		return
	}

	http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), file)

	return
}

func (c *Context) HttpAttachment(file http.File, name string) error {
	return c.contentDisposition(file, name, gox.ContentDispositionTypeAttachment)
}

func (c *Context) HttpInline(file http.File, name string) error {
	return c.contentDisposition(file, name, gox.ContentDispositionTypeInline)
}

func (c *Context) RequestBodyString() (body string, err error) {
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(c.Request().Body); nil != err {
		return
	}
	body = buf.String()

	return
}

func (c *Context) RequestBodyBytes() (body []byte, err error) {
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(c.Request().Body); nil != err {
		return
	}
	body = buf.Bytes()

	return
}

func (c *Context) contentDisposition(file http.File, name string, dispositionType gox.ContentDispositionType) error {
	c.Response().Header().Set(gox.HeaderContentDisposition, gox.ContentDisposition(name, dispositionType))

	return c.HttpFile(file)
}
