package echox

import (
	`fmt`
)

var _ option = (*optionAddr)(nil)

type optionAddr struct {
	// 主机
	host string
	// 端口
	port int
}

// RandomAddr 随机端口
func RandomAddr() *optionAddr {
	return Addr("", 0)
}

// Addr 绑定地址
func Addr(host string, port int) *optionAddr {
	return &optionAddr{
		host: host,
		port: port,
	}
}

func (a *optionAddr) apply(options *options) {
	options.addr = fmt.Sprintf("%s:%d", a.host, a.port)
}
