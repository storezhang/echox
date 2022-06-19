package echox

import (
	"fmt"
	"strings"

	"github.com/goexl/gox"
)

var (
	_ option = (*optionProxy)(nil)
	_        = Proxy
	_        = HttpProxy
	_        = HttpsProxy
	_        = ProxyAddr
)

type optionProxy struct {
	proxy string
}

// Proxy 配置代理
func Proxy(scheme gox.URIScheme, domain string, port int) *optionProxy {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%s://%s`, scheme, domain))
	// 处理默认端口
	if gox.URISchemeHttps == scheme && 443 != port || gox.URISchemeHttp == scheme && 80 != port {
		sb.WriteString(fmt.Sprintf(`:%d`, port))
	}

	return &optionProxy{
		proxy: sb.String(),
	}
}

// ProxyAddr 配置代理
func ProxyAddr(addr string) *optionProxy {
	return &optionProxy{
		proxy: addr,
	}
}

// HttpProxy 配置Http代理
func HttpProxy(domain string) *optionProxy {
	return &optionProxy{
		proxy: fmt.Sprintf(`http://%s`, domain),
	}
}

// HttpsProxy 配置Https代理
func HttpsProxy(domain string) *optionProxy {
	return &optionProxy{
		proxy: fmt.Sprintf(`https://%s`, domain),
	}
}

func (p *optionProxy) apply(options *options) {
	options.proxy = p.proxy
}
