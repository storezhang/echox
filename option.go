package echox

var defaultOptions = &options{
	addr: `:9000`,
	binder: &binder{
		tagParam:   `param`,
		tagQuery:   `query`,
		tagForm:    `form`,
		tagHeader:  `header`,
		tagDefault: `default`,
	},
	validate:   true,
	banner:     false,
	error:      errorHandlerFunc,
	crosEnable: true,
	cros: crosConfig{
		origins:     []string{`*`},
		credentials: true,
	},
	panicStack: panicStackConfig{
		size:              4 << 10,
		disableStackAll:   false,
		disablePrintStack: false,
	},
}

type (
	option interface {
		apply(options *options)
	}

	options struct {
		// 地址
		addr string
		// 代理
		proxy string
		// 是否需要使用绑定功能
		binder *binder
		// 是否需要使用数据验证功能
		validate bool
		// 是否显示徽标
		banner bool
		// 上下文路径，同一个机器上同一个端口运行不同的应用时需要加一个前级来区分不同的应用
		context string
		// 初始化方法
		inits []initFunc
		// 错误处理
		error errorHandler

		// 跨域
		crosEnable bool
		cros       crosConfig

		// 各种中间件
		// 打印堆栈信息
		panicStack panicStackConfig
	}
)

// NewOptions 创建空选项列表
func NewOptions(opts ...option) []option {
	return opts
}
