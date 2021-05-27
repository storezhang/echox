package echox

type restfulHandler func(ctx *Context) (rsp interface{}, err error)
