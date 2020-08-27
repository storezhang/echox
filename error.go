package echox

import (
	"github.com/storezhang/gox"
)

var (
	ErrNoUpdateParam = &gox.CodeError{ErrorCode: 100001, Msg: "未下发更新参数"}
)
