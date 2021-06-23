package echox

import (
	`time`
)

var _ startOption = (*startOptionShutdownTimeout)(nil)

type startOptionShutdownTimeout struct {
	// 超时
	timeout time.Duration
}

// ShutdownTimeout 配置退出超时时间
func ShutdownTimeout(timeout time.Duration) *startOptionShutdownTimeout {
	return &startOptionShutdownTimeout{
		timeout: timeout,
	}
}

func (st *startOptionShutdownTimeout) applyStart(options *startOptions) {
	options.shutdownTimeout = st.timeout
}
