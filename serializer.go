package echox

const (
	serializerJson    serializer = "json"
	serializerProto   serializer = "proto"
	serializerMsgpack serializer = "msgpack"
	serializerXml     serializer = "xml"
	serializerBytes   serializer = "bytes"
)

type serializer string
