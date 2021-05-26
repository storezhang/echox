package echox

var _ option = (*optionValidate)(nil)

type optionValidate struct{}

// DisableValidate 禁用数据验证
func DisableValidate() *optionValidate {
	return &optionValidate{}
}

func (t *optionValidate) apply(options *options) {
	options.validate = false
}
