package echox

import (
	`fmt`

	`github.com/labstack/echo/v4/middleware`
	`github.com/storezhang/gox`
)

var (
	// 支持的算法
	supportAlgorithm = []Algorithm{
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

	_ option = (*optionSignature)(nil)
)

type optionSignature struct {
	//  确定是不是要走中间件
	skipper middleware.Skipper
	//  签名算法
	algorithm Algorithm
	//  获得签名参数
	source keySource
}

// Signature Http签名
func Signature(algorithm Algorithm, source keySource) *optionSignature {
	return SignatureWithConfig(middleware.DefaultSkipper, algorithm, source)
}

// SignatureWithConfig Http签名
func SignatureWithConfig(skipper middleware.Skipper, algorithm Algorithm, source keySource) *optionSignature {
	// 检查算法配置是否正确
	if exists, _ := gox.IsInArray(algorithm, supportAlgorithm); !exists {
		panic(fmt.Errorf("不支持的算法：%s", algorithm))
	}

	return &optionSignature{
		skipper:   skipper,
		algorithm: algorithm,
		source:    source,
	}
}

func (j *optionSignature) apply(options *options) {
	options.signature.skipper = j.skipper
	options.signature.algorithm = j.algorithm
	options.signature.source = j.source
	options.signatureEnable = true
}
