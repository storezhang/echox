package echox

import (
	`fmt`

	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
)

// 支持的算法
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

type Signature struct {
	//  确定是不是要走中间件
	skipper middleware.Skipper `validate:"required"`
	//  签名算法
	algorithm Algorithm `validate:"required"`
	//  获得签名参数
	source keySource `validate:"required"`
}

// NewSignature Http签名
func NewSignature(algorithm Algorithm, source keySource) *Signature {
	return NewSignatureWithConfig(middleware.DefaultSkipper, algorithm, source)
}

// NewSignatureWithConfig Http签名
func NewSignatureWithConfig(skipper middleware.Skipper, algorithm Algorithm, source keySource) *Signature {
	// 检查算法配置是否正确
	if exists, _ := gox.IsInArray(algorithm, supportAlgorithm); !exists {
		panic(fmt.Errorf("不支持的算法：%s", algorithm))
	}

	return &Signature{
		skipper:   skipper,
		algorithm: algorithm,
		source:    source,
	}
}
