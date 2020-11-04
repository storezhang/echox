package echox

import (
	`crypto`
	`encoding/json`
	`fmt`
	`net/http`

	`github.com/go-fed/httpsig`
	`github.com/labstack/echo/v4`
	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
)

const (
	// 对称加密算法
	HmacWithSHA224     Algorithm = "hmac-sha224"
	HmacWithSHA256     Algorithm = "hmac-sha256"
	HmacWithSHA384     Algorithm = "hmac-sha384"
	HmacWithSHA512     Algorithm = "hmac-sha512"
	HmacWithRipemd160  Algorithm = "hmac-ripemd160"
	HmacWithSHA3224    Algorithm = "hmac-sha3-224"
	HmacWithSHA3256    Algorithm = "hmac-sha3-256"
	HmacWithSHA3384    Algorithm = "hmac-sha3-384"
	HmacWithSHA3512    Algorithm = "hmac-sha3-512"
	HmacWithSHA512224  Algorithm = "hmac-sha512-224"
	HmacWithSHA512256  Algorithm = "hmac-sha512-256"
	HmacWithBlake2s256 Algorithm = "hmac-blake2s-256"
	HmacWithBlake2b256 Algorithm = "hmac-blake2b-256"
	HmacWithBlake2b384 Algorithm = "hmac-blake2b-384"
	HmacWithBlake2b512 Algorithm = "hmac-blake2b-512"
	Blake2sWith256     Algorithm = "blake2s-256"
	Blake2bWith256     Algorithm = "blake2b-256"
	Blake2bWith384     Algorithm = "blake2b-384"
	Blake2bWith512     Algorithm = "blake2b-512"
	// RAS非对称加密算法
	RsaWithSHA224    Algorithm = "rsa-sha224"
	RsaWithSHA256    Algorithm = "rsa-sha256"
	RsaWithSHA384    Algorithm = "rsa-sha384"
	RsaWithSHA512    Algorithm = "rsa-sha512"
	RsaWithRipemd160 Algorithm = "rsa-ripemd160"
	// ECDSA非对称加密算法
	EcdsaWithSHA224    Algorithm = "ecdsa-sha224"
	EcdsaWithSHA256    Algorithm = "ecdsa-sha256"
	EcdsaWithSHA384    Algorithm = "ecdsa-sha384"
	EcdsaWithSHA512    Algorithm = "ecdsa-sha512"
	EcdsaWithRipemd160 Algorithm = "ecdsa-ripemd160"
)

var supportAlgorithm = []Algorithm{
	// 对称加密算法
	HmacWithSHA224,
	HmacWithSHA256,
	HmacWithSHA384,
	HmacWithSHA512,
	HmacWithRipemd160,
	HmacWithSHA3224,
	HmacWithSHA3256,
	HmacWithSHA3384,
	HmacWithSHA3512,
	HmacWithSHA512224,
	HmacWithSHA512256,
	HmacWithBlake2s256,
	HmacWithBlake2b256,
	HmacWithBlake2b384,
	HmacWithBlake2b512,
	Blake2sWith256,
	Blake2bWith256,
	Blake2bWith384,
	Blake2bWith512,
	// RAS非对称加密算法
	RsaWithSHA224,
	RsaWithSHA256,
	RsaWithSHA384,
	RsaWithSHA512,
	RsaWithRipemd160,
	// ECDSA非对称加密算法
	EcdsaWithSHA224,
	EcdsaWithSHA256,
	EcdsaWithSHA384,
	EcdsaWithSHA512,
	EcdsaWithRipemd160,
}

var DefaultSignatureConfig = SignatureConfig{
	Skipper:   middleware.DefaultSkipper,
	Algorithm: HmacWithSHA512,
}

type (
	// SignatureConfig 签名配置
	SignatureConfig struct {
		// Skipper 确定是不是要走中间件
		Skipper middleware.Skipper `json:"skipper"`
		// Algorithm 签名算法
		Algorithm Algorithm `json:"algorithm"`
		// KeySource 获得签名参数
		KeySource KeySource `json:"key_source"`
	}

	// KeySource 获得签名参数
	KeySource interface {
		// SecretKey 获得签名参数
		SecretKey(id string) (secretKey string, err error)
	}

	// Algorithm 签名算法
	Algorithm string
)

// SignatureMiddleware 默认鉴权
func SignatureMiddleware() echo.MiddlewareFunc {
	return SignatureWithConfig(DefaultSignatureConfig)
}

// SignatureWithConfig 按配置鉴权
func SignatureWithConfig(config SignatureConfig) echo.MiddlewareFunc {
	if nil == config.Skipper {
		config.Skipper = DefaultSignatureConfig.Skipper
	}

	// 检查算法配置是否正确
	if exists, _ := gox.IsInArray(config.Algorithm, supportAlgorithm); !exists {
		panic(fmt.Errorf("不支持的算法：%s", config.Algorithm))
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if config.Skipper(ctx) {
				return next(ctx)
			}

			req := ctx.Request()

			var (
				verifier  httpsig.Verifier
				secretKey string
			)
			if verifier, err = httpsig.NewVerifier(req); nil != err {
				return
			}

			appKey := verifier.KeyId()
			if secretKey, err = config.KeySource.SecretKey(appKey); nil != err {
				return
			}

			key := crypto.PublicKey([]byte(secretKey))
			algorithm := httpsig.Algorithm(config.Algorithm)
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

func (sc SignatureConfig) String() string {
	jsonBytes, _ := json.MarshalIndent(sc, "", "    ")

	return string(jsonBytes)
}
