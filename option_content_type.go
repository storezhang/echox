package echox

var (
	_ httpOption  = (*optionContentType)(nil)
	_ groupOption = (*optionContentType)(nil)
)

type optionContentType struct {
	contentType string
}

// ContentType 配置类型
func ContentType(contentType string) *optionContentType {
	return &optionContentType{
		contentType: contentType,
	}
}

func (s *optionContentType) applyHttp(options *httpOptions) {
	options.contentType = s.contentType
}

func (s *optionContentType) applyGroup(options *groupOptions) {
	options.contentType = s.contentType
}
