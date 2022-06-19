package echox

import (
	"bytes"
	"net/http"
	"os"
	"strconv"

	"github.com/goexl/gox"
	"github.com/labstack/echo/v4"
)

// Context 自定义的Echo上下文
type Context struct {
	echo.Context
}

func (c *Context) IntParam(name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}

func (c *Context) Int64Param(name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}

func (c *Context) Data(rsp interface{}, opts ...httpOption) error {
	_options := defaultHttpOptions()
	for _, opt := range opts {
		opt.applyHttp(_options)
	}

	return data(c.Context, rsp, _options)
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

func (c *Context) BodyString() (body string, err error) {
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(c.Request().Body); nil != err {
		return
	}
	body = buf.String()

	return
}

func (c *Context) BodyBytes() (body []byte, err error) {
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

func parseContext(ctx echo.Context) (context *Context) {
	if _ctx, ok := ctx.(*Context); ok {
		context = _ctx
	} else {
		context = &Context{Context: ctx}
	}

	return
}
