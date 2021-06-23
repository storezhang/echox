package echox

import (
	`fmt`
	`net/http`
	`reflect`
	`strings`

	`github.com/dgrijalva/jwt-go`
	`github.com/labstack/echo/v4`
)

const (
	// AlgorithmHS256 HS256加密算法
	AlgorithmHS256 = "HS256"
)

var errJwtMissing = echo.NewHTTPError(http.StatusUnauthorized, "缺失Jwt请求头")

type (
	// 成功后的处理方法
	jwtSuccessHandler func(echo.Context)
	jwtExtractor      func(echo.Context) (string, error)
)

// JwtMiddleware Jwt中间件
func JwtMiddleware(config Jwt) echo.MiddlewareFunc {
	config.keyFunc = func(t *jwt.Token) (key interface{}, err error) {
		if t.Method.Alg() != config.method {
			err = fmt.Errorf("未知的签名算法=%v", t.Header["alg"])
		} else {
			key = []byte(config.key.(string))
		}

		return
	}

	for _, lookup := range config.lookups {
		parts := strings.Split(lookup, ":")
		switch parts[0] {
		case "header":
			config.extractor = append(config.extractor, jwtFromHeader(parts[1], config.scheme))
		case "query":
			config.extractor = append(config.extractor, jwtFromQuery(parts[1]))
		case "cookie":
			config.extractor = append(config.extractor, jwtFromCookie(parts[1]))
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if config.skipper(ctx) {
				err = next(ctx)

				return
			}

			if nil != config.beforeHandler {
				config.beforeHandler(ctx)
			}

			var authToken string
			if authToken, err = config.runExtractor(ctx); nil != err {
				return
			}

			token := new(jwt.Token)
			if _, ok := config.claims.(jwt.MapClaims); ok {
				token, err = jwt.Parse(authToken, config.keyFunc)
			} else {
				t := reflect.ValueOf(config.claims).Type().Elem()
				claims := reflect.New(t).Interface().(jwt.Claims)
				token, err = jwt.ParseWithClaims(authToken, claims, config.keyFunc)
			}
			if nil != err {
				return
			}

			if token.Valid {
				ctx.Set(config.context, token)
				if nil != config.successHandler {
					config.successHandler(ctx)
				}
				err = next(ctx)
			} else {
				err = &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "JWT错误或者已经失效",
					Internal: err,
				}
			}

			return
		}
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(ctx echo.Context) (token string, err error) {
		auth := ctx.Request().Header.Get(header)
		authLength := len(authScheme)
		if len(auth) > authLength+1 && auth[:authLength] == authScheme {
			token = auth[authLength+1:]
		} else {
			err = errJwtMissing
		}

		return
	}
}

func jwtFromQuery(param string) jwtExtractor {
	return func(ctx echo.Context) (token string, err error) {
		token = ctx.QueryParam(param)
		if "" == token {
			err = errJwtMissing
		}

		return
	}
}

func jwtFromCookie(name string) jwtExtractor {
	return func(ctx echo.Context) (token string, err error) {
		var cookie *http.Cookie
		if cookie, err = ctx.Cookie(name); nil != err {
			err = errJwtMissing
		} else {
			token = cookie.Value
		}

		return
	}
}
