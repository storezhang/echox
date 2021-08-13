package echox

var _ startOption = (*startOptionGraceful)(nil)

type startOptionGraceful struct{}

// Graceful 配置Graceful退出机制
func Graceful() *startOptionGraceful {
	return &startOptionGraceful{}
}

func (g *startOptionGraceful) applyStart(startOptions *startOptions) {
	startOptions.graceful = true
}
