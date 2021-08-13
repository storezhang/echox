package echox

type stopOption interface {
	applyStop(options *stopOptions)
}
