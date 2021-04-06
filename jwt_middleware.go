package echox

import (
	`encoding/json`
	`fmt`
	`net/http`
	`reflect`
	`strings`
	`time`

	`github.com/dgrijalva/jwt-go`
	`github.com/labstack/echo/v4`
	`github.com/labstack/echo/v4/middleware`
	`github.com/rs/xid`
)

const (
	// AlgorithmHS256 HS256加密算法
	AlgorithmHS256 = "HS256"
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "缺失JWT请求头")

	// DefaultJWTConfig 默认配置
	DefaultJWTConfig = &JWTConfig{
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   []string{"header:" + echo.HeaderAuthorization, "query:token"},
		AuthScheme:    "Bearer",
		Claims:        &jwt.StandardClaims{},
	}
)

type (
	// JWTConfig JWT中间件的配置
	JWTConfig struct {
		// 确定是不是要走中间件
		Skipper middleware.Skipper

		// 执行前的操作
		BeforeFunc middleware.BeforeFunc

		// 成功后操作
		SuccessHandler JWTSuccessHandler

		// 错误处理
		// 一般拿来返回自定义的JSON格式
		ErrorHandler JWTErrorHandler

		// 签名密钥
		// 必须字段
		SigningKey interface{}

		// 签名方法
		// 非必须 默认是HS256
		SigningMethod string

		// 存储用户信息的键
		// 非必须 默认值是"user"
		ContextKey string

		// 存储数据的类型
		// 非必须 默认值是jwt.MapClaims
		Claims jwt.Claims

		// 定义从哪获得Token
		// 非必须 默认值是"header:Authorization"
		// 可能的值：
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup []string

		// Token分隔字符串
		// 非必须 默认值是"Bearer"
		AuthScheme string

		keyFunc jwt.Keyfunc

		extractor []jwtExtractor
	}

	// JWTSuccessHandler 成功后的处理
	JWTSuccessHandler func(echo.Context)

	// JWTErrorHandler 错误处理
	JWTErrorHandler func(error) error

	jwtExtractor func(echo.Context) (string, error)
)

func (jc JWTConfig) String() string {
	jsonBytes, _ := json.MarshalIndent(jc, "", "    ")

	return string(jsonBytes)
}

func (jc *JWTConfig) Parse(t string) (claims jwt.Claims, header map[string]interface{}, err error) {
	token := new(jwt.Token)
	if _, ok := jc.Claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(t, jc.keyFunc)
	} else {
		elem := reflect.ValueOf(jc.Claims).Type().Elem()
		claims := reflect.New(elem).Interface().(jwt.Claims)
		token, err = jwt.ParseWithClaims(t, claims, jc.keyFunc)
	}
	if nil == err && token.Valid {
		claims = token.Claims
		header = token.Header
	}

	return
}

func (jc *JWTConfig) Extractor(c echo.Context) (token string, err error) {
	for _, extractor := range jc.extractor {
		if token, err = extractor(c); nil == err || "" != token {
			break
		}
	}

	return
}

func (jc *JWTConfig) Token(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(jc.SigningMethod), claims)

	return token.SignedString([]byte(jc.SigningKey.(string)))
}

func (jc *JWTConfig) UserToken(
	domain string,
	data interface{},
	expire time.Duration,
) (token string, id string, err error) {
	// 序列化User对象为JSON
	var userBytes []byte
	if userBytes, err = json.Marshal(data); nil != err {
		return
	}

	id = xid.New().String()
	token, err = jc.Token(jwt.StandardClaims{
		// 代表这个JWT的签发主体
		Issuer: domain,
		// 代表这个JWT的主体，即它的所有人
		Subject: string(userBytes),
		// 代表这个JWT的接收对象
		Audience: domain,
		// 是一个时间戳，代表这个JWT的签发时间
		IssuedAt: time.Now().Unix(),
		// 是一个时间戳，代表这个JWT生效的开始时间，意味着在这个时间之前验证JWT是会失败的
		NotBefore: time.Now().Unix(),
		// 是一个时间戳，代表这个JWT的过期时间
		ExpiresAt: time.Now().Add(expire).Unix(),
		// 是JWT的唯一标识
		Id: id,
	})

	return
}

// JWT JWT中间件
func JWT(key interface{}) echo.MiddlewareFunc {
	c := DefaultJWTConfig
	c.SigningKey = key

	return JWTWithConfig(c)
}

// JWTWithConfig JWT中间件
func JWTWithConfig(config *JWTConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil {
		panic("echo: jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if 0 == len(config.TokenLookup) {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(config.SigningKey.(string)), nil
	}

	for _, tokenLookup := range config.TokenLookup {
		parts := strings.Split(tokenLookup, ":")
		switch parts[0] {
		case "header":
			config.extractor = append(config.extractor, jwtFromHeader(parts[1], config.AuthScheme))
		case "query":
			config.extractor = append(config.extractor, jwtFromQuery(parts[1]))
		case "cookie":
			config.extractor = append(config.extractor, jwtFromCookie(parts[1]))
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			auth, err := config.Extractor(c)
			if err != nil {
				if config.ErrorHandler != nil {
					return config.ErrorHandler(err)
				}
				return err
			}
			token := new(jwt.Token)
			// Issue #647, #656
			if _, ok := config.Claims.(jwt.MapClaims); ok {
				token, err = jwt.Parse(auth, config.keyFunc)
			} else {
				t := reflect.ValueOf(config.Claims).Type().Elem()
				claims := reflect.New(t).Interface().(jwt.Claims)
				token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
			}
			if err == nil && token.Valid {
				c.Set(config.ContextKey, token)
				if config.SuccessHandler != nil {
					config.SuccessHandler(c)
				}

				return next(c)
			}
			if config.ErrorHandler != nil {
				return config.ErrorHandler(err)
			}

			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "JWT错误或者已经失效",
				Internal: err,
			}
		}
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}

		return "", ErrJWTMissing
	}
}

func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		if token == "" {
			return "", ErrJWTMissing
		}

		return token, nil
	}
}

func jwtFromCookie(name string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", ErrJWTMissing
		}
		return cookie.Value, nil
	}
}
