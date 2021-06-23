package echox

type startOption interface {
	applyStart(options *startOptions)
}
