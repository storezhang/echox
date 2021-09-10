package echox

type (
	groupOption interface {
		applyGroup(options *groupOptions)
	}

	groupOptions struct {
		*httpOptions

		middlewares []MiddlewareFunc
	}
)

func defaultGroupOptions() *groupOptions {
	return &groupOptions{
		httpOptions: defaultHttpOptions(),

		middlewares: make([]MiddlewareFunc, 0, 0),
	}
}
