package echox

import (
	`encoding/json`
	`io/ioutil`
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
		JWT *JWTConfig
	}
)

func (ec *EchoContext) User() (user gox.BaseUser, err error) {
	var (
		token  string
		claims jwt.Claims
	)

	if token, err = ec.JWT.Extractor(ec.Context); nil != err {
		return
	}

	if claims, _, err = ec.JWT.Parse(token); nil != err {
		return
	}

	// 从JWT Token中反序列化User
	err = json.Unmarshal([]byte(claims.(*jwt.StandardClaims).Subject), &user)

	return
}

func (ec *EchoContext) JWTToken(domain string, user gox.BaseUser) (token string, err error) {
	// 序列化User对象为JSON
	var userBytes []byte
	if userBytes, err = json.Marshal(user); nil != err {
		return
	}

	if token, err = ec.JWT.Token(jwt.StandardClaims{
		// 代表这个JWT的签发主体
		Issuer: domain,
		// 代表这个JWT的主体，即它的所有人
		Subject: string(userBytes),
		// 代表这个JWT的接收对象
		Audience: domain,
		// 是一个时间戳，代表这个JWT的过期时间
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}); nil != err {
		return
	}

	return
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
	if indent != "" {
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

// 获取有关联表的更新信息
func UpdateWithRelation(c echo.Context, bean interface{}, notCols ...string) (cols, otherCols []string, err error) {
	var (
		reqMap = make(map[string]interface{})
	)

	if err = UpdateMap(c, bean, &reqMap); nil != err {
		return
	}

	cols = make([]string, 0)
	otherCols = make([]string, 0)
	for key := range reqMap {
		if exists, _ := gox.IsInArray(key, notCols); exists {
			otherCols = append(otherCols, gox.UnderscoreName(key, false))
		} else {
			cols = append(cols, gox.UnderscoreName(key, false))
		}
	}

	if 0 == len(cols) && 0 == len(otherCols) {
		err = ErrNoUpdateParam
	}

	return
}

func UpdateInfo(c echo.Context, bean interface{}) (cols []string, err error) {
	var reqMap = make(map[string]interface{})

	if err = UpdateMap(c, bean, &reqMap); nil != err {
		return
	}

	cols = make([]string, 0)
	for key := range reqMap {
		cols = append(cols, gox.UnderscoreName(key, false))
	}

	if 0 == len(cols) {
		err = ErrNoUpdateParam
	}

	return
}

func UpdateMap(c echo.Context, bean, reqMap interface{}) (err error) {
	var body []byte

	if body, err = ioutil.ReadAll(c.Request().Body); nil != err {
		return
	}
	if err = json.Unmarshal(body, bean); nil != err {
		return
	}
	if err = json.Unmarshal(body, &reqMap); nil != err {
		return
	}

	return
}
