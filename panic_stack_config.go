package echox

type panicStackConfig struct {
	// 方法栈大小
	// 默认4KB
	size int `validate:"required"`
	// 是否禁止显示所有的栈信息
	disableStackAll bool `validate:"required"`
	// 禁止打印栈信息
	disablePrintStack bool `validate:"required"`
}
