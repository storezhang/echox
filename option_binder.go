package echox

var _ option = (*optionBinder)(nil)

type optionBinder struct {
	param   string
	query   string
	form    string
	disable bool
}

// Binder 配置数据绑定
func Binder(param string, query string, form string) *optionBinder {
	return &optionBinder{
		param: param,
		query: query,
		form:  form,
	}
}

// DisableBinder 禁用数据绑定
func DisableBinder() *optionBinder {
	return &optionBinder{
		disable: true,
	}
}

func (b *optionBinder) apply(options *options) {
	if b.disable {
		options.binder = nil
	} else {
		options.binder.tagParam = b.param
		options.binder.tagQuery = b.query
		options.binder.tagForm = b.form
	}
}
