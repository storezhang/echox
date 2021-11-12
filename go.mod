module github.com/storezhang/echox/v2

go 1.16

require (
	github.com/casbin/casbin/v2 v2.30.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-fed/httpsig v1.1.0
	github.com/go-playground/validator/v10 v10.6.1
	github.com/labstack/echo/v4 v4.3.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/rs/xid v1.3.0
	github.com/storezhang/god v0.0.2
	github.com/storezhang/gox v1.8.1
	github.com/storezhang/validatorx v1.0.8
	github.com/vmihailenco/msgpack/v5 v5.3.4
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/protobuf v1.27.1
)

// replace github.com/storezhang/validatorx => ../../storezhang/validatorx
// replace github.com/storezhang/gox => ../../storezhang/gox
