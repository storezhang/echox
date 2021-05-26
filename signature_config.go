package echox

import (
	`github.com/labstack/echo/v4/middleware`
)

type signatureConfig struct {
	//  确定是不是要走中间件
	skipper middleware.Skipper `validate:"required"`
	//  签名算法
	algorithm Algorithm `validate:"required"`
	//  获得签名参数
	source keySource `validate:"required"`
}
