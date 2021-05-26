package echox

import (
	`time`
)

var _ option = (*optionShutdownTimeout)(nil)

type optionShutdownTimeout struct {
	// 超时
	timeout time.Duration
}

// ShutdownTimeout 配置退出超时时间
func ShutdownTimeout(timeout time.Duration) *optionShutdownTimeout {
	return &optionShutdownTimeout{
		timeout: timeout,
	}
}

func (st *optionShutdownTimeout) apply(options *options) {
	options.shutdownTimeout = st.timeout
}
