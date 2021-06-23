package echox

import (
	`crypto`
	`net/http`

	`github.com/go-fed/httpsig`
	`github.com/labstack/echo/v4`
)

// keySource 获得签名参数
type keySource interface {
	// Key 获得签名参数
	Key(id string) (key string, err error)
}

// SignatureMiddleware 签名中间件
func SignatureMiddleware(config Signature) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if config.skipper(ctx) {
				err = next(ctx)

				return
			}

			req := ctx.Request()

			var verifier httpsig.Verifier
			if verifier, err = httpsig.NewVerifier(req); nil != err {
				return
			}

			appKey := verifier.KeyId()
			var secretKey string
			if secretKey, err = config.source.Key(appKey); nil != err {
				return
			}

			key := crypto.PublicKey([]byte(secretKey))
			algorithm := httpsig.Algorithm(config.algorithm)
			if err = verifier.Verify(key, algorithm); nil != err {
				err = &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "未经允许，禁止驶入！",
					Internal: err,
				}
			} else {
				err = next(ctx)
			}

			return
		}
	}
}
