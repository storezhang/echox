package echox

var (
	_ httpOption  = (*optionSerializer)(nil)
	_ groupOption = (*optionSerializer)(nil)
)

type optionSerializer struct {
	serializer serializer
}

// Proto 谷歌Protocol Buffer序列化
func Proto() *optionSerializer {
	return &optionSerializer{
		serializer: serializerProto,
	}
}

// JSON 使用JSON序列化
func JSON() *optionSerializer {
	return &optionSerializer{
		serializer: serializerJson,
	}
}

// XML 使用XML序列化
func XML() *optionSerializer {
	return &optionSerializer{
		serializer: serializerXml,
	}
}

// Msgpack 使用Msgpack序列化
func Msgpack() *optionSerializer {
	return &optionSerializer{
		serializer: serializerMsgpack,
	}
}

// Bytes 原始数据
func Bytes() *optionSerializer {
	return &optionSerializer{
		serializer: serializerBytes,
	}
}

func (s *optionSerializer) applyHttp(options *httpOptions) {
	options.serializer = s.serializer
}

func (s *optionSerializer) applyGroup(options *groupOptions) {
	options.serializer = s.serializer
}
