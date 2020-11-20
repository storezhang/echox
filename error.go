package echox

import (
	"github.com/storezhang/gox"
)

var (
	// ErrNoUpdateParam 未下发更新参数
	ErrNoUpdateParam = &gox.CodeError{ErrorCode: 100001, Message: "未下发更新参数"}
)
