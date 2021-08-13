package echox

import (
	`time`
)

type stopOptions struct {
	*options

	// 退出超时时间
	timeout time.Duration
}

func defaultStopOptions() *stopOptions {
	return &stopOptions{
		options: defaultOptions,

		timeout: 30 * time.Second,
	}
}
