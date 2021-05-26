package echox

var _ option = (*optionBinder)(nil)

type optionBinder struct{}

// DisableBinder 禁用数据绑定
func DisableBinder() *optionBinder {
	return &optionBinder{}
}

func (t *optionBinder) apply(options *options) {
	options.binder = false
}
