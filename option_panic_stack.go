package echox

var _ option = (*optionPanicStack)(nil)

type optionPanicStack struct {
	// 方法栈大小
	// 默认4KB
	size int
	// 是否禁止显示所有的栈信息
	disableStackAll bool
	// 禁止打印栈信息
	disablePrintStack bool
}

// PanicStack 堆栈打印快捷方式
func PanicStack(size int) *optionPanicStack {
	return PanicStackWithConfig(size, false, false)
}

// PanicStackWithConfig 堆栈打印
func PanicStackWithConfig(size int, disableStackAll bool, disablePrintStack bool) *optionPanicStack {
	return &optionPanicStack{
		size:              size,
		disableStackAll:   disableStackAll,
		disablePrintStack: disablePrintStack,
	}
}

func (ps *optionPanicStack) apply(options *options) {
	options.panicStack.size = ps.size
	options.panicStack.disableStackAll = ps.disableStackAll
	options.panicStack.disablePrintStack = ps.disablePrintStack
}
