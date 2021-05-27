package echox

import (
	`encoding/json`
	`net/http`
	`os`
	`time`

	`github.com/dgrijalva/jwt-go`
	`github.com/json-iterator/go`
	`github.com/labstack/echo/v4`
	`github.com/storezhang/gox`
)

const defaultIndent = "  "

// Context 自定义的Echo上下文
type Context struct {
	echo.Context

	// Jwt配置
	jwt JwtConfig
}

func (c *Context) Subject(subject interface{}) (err error) {
	var (
		token  string
		claims jwt.Claims
	)

	if token, err = c.jwt.runExtractor(c.Context); nil != err {
		return
	}
	if claims, _, err = c.jwt.Parse(token); nil != err {
		return
	}
	// 从Token中反序列化主题数据
	err = json.Unmarshal([]byte(claims.(*jwt.StandardClaims).Subject), &subject)

	return
}

func (c *Context) JwtToken(domain string, data interface{}, expire time.Duration) (token string, id string, err error) {
	return c.jwt.Token(domain, data, expire)
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

func (c *Context) contentDisposition(file http.File, name string, dispositionType gox.ContentDispositionType) error {
	c.Response().Header().Set(gox.HeaderContentDisposition, gox.ContentDisposition(name, dispositionType))

	return c.HttpFile(file)
}

func (c *Context) JSON(code int, data interface{}) error {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.Echo().Debug || pretty {
		indent = defaultIndent
	}

	return c.json(code, data, indent)
}

func (c *Context) JSONPretty(code int, data interface{}, indent string) (err error) {
	return c.json(code, data, indent)
}

func (c *Context) JSONBlob(code int, data []byte) (err error) {
	return c.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, data)
}

func (c *Context) JSONP(code int, callback string, data interface{}) (err error) {
	return c.jsonPBlob(code, callback, data)
}

func (c *Context) JSONPBlob(code int, callback string, data []byte) (err error) {
	c.writeContentType(echo.MIMEApplicationJavaScriptCharsetUTF8)
	c.Response().WriteHeader(code)
	if _, err = c.Response().Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.Response().Write(data); err != nil {
		return
	}
	_, err = c.Response().Write([]byte(");"))

	return
}

func (c *Context) jsonPBlob(code int, callback string, data interface{}) (err error) {
	enc := jsoniter.NewEncoder(c.Response())
	_, pretty := c.QueryParams()["pretty"]
	if c.Echo().Debug || pretty {
		enc.SetIndent("", "  ")
	}
	c.writeContentType(echo.MIMEApplicationJavaScriptCharsetUTF8)
	c.Response().WriteHeader(code)
	if _, err = c.Response().Write([]byte(callback + "(")); err != nil {
		return
	}
	if err = enc.Encode(data); err != nil {
		return
	}
	if _, err = c.Response().Write([]byte(");")); err != nil {
		return
	}

	return
}

func (c *Context) json(code int, data interface{}, indent string) error {
	enc := jsoniter.NewEncoder(c.Response())
	if "" != indent {
		enc.SetIndent("", indent)
	}
	c.writeContentType(echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(code)

	return enc.Encode(data)
}

func (c *Context) writeContentType(value string) {
	header := c.Response().Header()
	if "" == header.Get(echo.HeaderContentType) {
		header.Set(echo.HeaderContentType, value)
	}
}
