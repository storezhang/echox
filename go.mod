module github.com/storezhang/echox

go 1.14

require (
	github.com/casbin/casbin/v2 v2.7.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-fed/httpsig v1.0.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/json-iterator/go v1.1.10
	github.com/labstack/echo/v4 v4.2.1
	github.com/mcuadros/go-defaults v1.2.0
	github.com/rs/xid v1.2.1
	github.com/storezhang/gox v1.3.9
	github.com/storezhang/validatorx v1.0.2
)

// replace github.com/storezhang/validatorx => ../../storezhang/validatorx
// replace github.com/storezhang/gox => ../../storezhang/gox
