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

const (
	defaultIndent = "  "
)

type (
	// EchoContext
	EchoContext struct {
		echo.Context

		// JWT配置
		jwt *JWTConfig
	}
)

func (ec *EchoContext) User() (user gox.BaseUser, err error) {
	var (
		token  string
		claims jwt.Claims
	)

	if token, err = ec.jwt.Extractor(ec.Context); nil != err {
		return
	}

	if claims, _, err = ec.jwt.Parse(token); nil != err {
		return
	}

	// 从JWT Token中反序列化User
	err = json.Unmarshal([]byte(claims.(*jwt.StandardClaims).Subject), &user)

	return
}

func (ec *EchoContext) JWTToken(domain string, user gox.BaseUser, expire time.Duration) (token string, id string, err error) {
	return ec.jwt.UserToken(domain, user, expire)
}

func (ec *EchoContext) HttpFile(file http.File) (err error) {
	defer func() {
		_ = file.Close()
	}()

	var fi os.FileInfo
	fi, err = file.Stat()
	if nil != err {
		return
	}

	http.ServeContent(ec.Response(), ec.Request(), fi.Name(), fi.ModTime(), file)

	return
}

func (ec *EchoContext) HttpAttachment(file http.File, name string) error {
	return ec.contentDisposition(file, name, gox.ContentDispositionTypeAttachment)
}

func (ec *EchoContext) HttpInline(file http.File, name string) error {
	return ec.contentDisposition(file, name, gox.ContentDispositionTypeInline)
}

func (ec *EchoContext) contentDisposition(file http.File, name string, dispositionType gox.ContentDispositionType) error {
	ec.Response().Header().Set(gox.HeaderContentDisposition, gox.ContentDisposition(name, dispositionType))

	return ec.HttpFile(file)
}

func (ec *EchoContext) JSON(code int, i interface{}) (err error) {
	indent := ""
	if _, pretty := ec.QueryParams()["pretty"]; ec.Echo().Debug || pretty {
		indent = defaultIndent
	}
	return ec.json(code, i, indent)
}

func (ec *EchoContext) JSONPretty(code int, i interface{}, indent string) (err error) {
	return ec.json(code, i, indent)
}

func (ec *EchoContext) JSONBlob(code int, b []byte) (err error) {
	return ec.Blob(code, echo.MIMEApplicationJSONCharsetUTF8, b)
}

func (ec *EchoContext) JSONP(code int, callback string, i interface{}) (err error) {
	return ec.jsonPBlob(code, callback, i)
}

func (ec *EchoContext) JSONPBlob(code int, callback string, b []byte) (err error) {
	ec.writeContentType(echo.MIMEApplicationJavaScriptCharsetUTF8)
	ec.Response().WriteHeader(code)
	if _, err = ec.Response().Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = ec.Response().Write(b); err != nil {
		return
	}
	_, err = ec.Response().Write([]byte(");"))

	return
}

func (ec *EchoContext) jsonPBlob(code int, callback string, i interface{}) (err error) {
	enc := jsoniter.NewEncoder(ec.Response())
	_, pretty := ec.QueryParams()["pretty"]
	if ec.Echo().Debug || pretty {
		enc.SetIndent("", "  ")
	}
	ec.writeContentType(echo.MIMEApplicationJavaScriptCharsetUTF8)
	ec.Response().WriteHeader(code)
	if _, err = ec.Response().Write([]byte(callback + "(")); err != nil {
		return
	}
	if err = enc.Encode(i); err != nil {
		return
	}
	if _, err = ec.Response().Write([]byte(");")); err != nil {
		return
	}

	return
}

func (ec *EchoContext) json(code int, i interface{}, indent string) error {
	enc := jsoniter.NewEncoder(ec.Response())
	if "" != indent {
		enc.SetIndent("", indent)
	}
	ec.writeContentType(echo.MIMEApplicationJSONCharsetUTF8)
	ec.Response().WriteHeader(code)

	return enc.Encode(i)
}

func (ec *EchoContext) writeContentType(value string) {
	header := ec.Response().Header()
	if "" == header.Get(echo.HeaderContentType) {
		header.Set(echo.HeaderContentType, value)
	}
}
