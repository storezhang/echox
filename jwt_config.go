package echox

// JwtConfig Jwt配置
type JwtConfig struct {
	// 签名密钥
	// 必须字段
	Key []byte `json:"key" yaml:"key" validate:"required"`
	// 签名方法
	// 非必须，默认是HS256
	Method string `default:"HS256" json:"method" yaml:"method" validate:"required"`
	// 存储用户信息的键
	// 非必须，默认值是"user"
	Context string `default:"user" json:"context" yaml:"context" validate:"required"`
	// 定义从哪获得Token
	// 非必须，默认值是"header:Authorization和query:token"
	// 可能的值：
	// "header:<name>"
	// "query:<name>"
	// "cookie:<name>"
	Lookups []string `default:"[header:Authorization,query:token]" json:"lookups" yaml:"lookups" validate:"required"`
	// Token分隔字符串
	// 非必须，默认值是"Bearer"
	Scheme string `default:"Bearer" json:"scheme" yaml:"scheme" validate:"required"`
}
