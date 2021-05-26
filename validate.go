package echox

import (
	`github.com/go-playground/validator/v10`
)

type validate struct {
	validate *validator.Validate
}

func (v *validate) Validate(data interface{}) error {
	return v.validate.Struct(data)
}
