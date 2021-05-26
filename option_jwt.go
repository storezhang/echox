package echox

var _ option = (*optionJwt)(nil)

type optionJwt struct {
	jwt JwtConfig
}

// Jwt 绑定地址
func Jwt(jwt JwtConfig) *optionJwt {
	return &optionJwt{
		jwt: jwt,
	}
}

func (j *optionJwt) apply(options *options) {
	options.jwt = j.jwt
	options.jwtEnable = true
}
