package echox

import (
	`time`
)

var (
	_ startOption = (*optionShutdownTimeout)(nil)
	_ stopOption  = (*optionShutdownTimeout)(nil)
)

type optionShutdownTimeout struct {
	timeout time.Duration
}

// ShutdownTimeout 配置退出超时时间
func ShutdownTimeout(timeout time.Duration) *optionShutdownTimeout {
	return &optionShutdownTimeout{
		timeout: timeout,
	}
}

func (st *optionShutdownTimeout) applyStart(options *startOptions) {
	options.shutdownTimeout = st.timeout
}

func (st *optionShutdownTimeout) applyStop(options *stopOptions) {
	options.timeout = st.timeout
}
