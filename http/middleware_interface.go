package http

type Middleware interface {
	// 如果返回needStop == true 则会结束路由调用
	HandlerBefore(*Context) (needStop bool)
	HandlerAfter(*Context)
}
