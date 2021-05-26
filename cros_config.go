package echox

type crosConfig struct {
	// 允许访问资源的地址
	origins []string
	// 是否可以将对请求的响应暴露给页面
	credentials bool
}
