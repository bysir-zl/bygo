package http


type Middleware interface {
    //如果返回needStop == true
    //则会结束路由调用,直接返回daa
    HandlerBefore(Context) (needStop bool,data ResponseData)
    HandlerAfter(Context) (needStop bool,data ResponseData)
}
