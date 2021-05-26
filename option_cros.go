package echox

var _ option = (*optionCros)(nil)

type optionCros struct {
	// 是否开启
	enable bool
	// 允许访问资源的地址
	origins []string
	// 是否可以将对请求的响应暴露给页面
	credentials bool
}

// Cros 跨域配置
func Cros(credentials bool, origins ...string) *optionCros {
	return &optionCros{
		enable:      true,
		origins:     origins,
		credentials: credentials,
	}
}

// DisableCros 关闭跨域
func DisableCros() *optionCros {
	return &optionCros{
		enable: false,
	}
}

func (c *optionCros) apply(options *options) {
	options.crosEnable = c.enable
	options.cros.origins = c.origins
	options.cros.credentials = c.credentials
}
