package echox

var _ option = (*optionDomain)(nil)

type optionDomain struct {
	domain string
}

// Domain 绑定地址
func Domain(domain string) *optionDomain {
	return &optionDomain{
		domain: domain,
	}
}

func (d *optionDomain) apply(options *options) {
	options.domain = d.domain
}
