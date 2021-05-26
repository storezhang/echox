package echox

var _ option = (*optionBanner)(nil)

type optionBanner struct{}

// Banner 开始徽标
func Banner() *optionBanner {
	return &optionBanner{}
}

func (t *optionBanner) apply(options *options) {
	options.banner = true
}
