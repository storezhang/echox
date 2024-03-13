module github.com/storezhang/echox/v2

go 1.16

require (
	github.com/casbin/casbin/v2 v2.47.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-fed/httpsig v1.1.0
	github.com/go-playground/validator/v10 v10.10.1
	github.com/goexl/exc v0.0.4
	github.com/goexl/gox v0.0.6
	github.com/goexl/mengpo v0.1.7
	github.com/goexl/xiren v0.0.3
	github.com/labstack/echo/v4 v4.9.0
	github.com/rs/xid v1.4.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/protobuf v1.33.0
)

// replace github.com/goexl/xiren => ../../storezhang/validatorx
// replace github.com/goexl/gox => ../../storezhang/gox
