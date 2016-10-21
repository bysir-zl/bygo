package router

import (
    "app/middleware"
    "app/controller"
    "bygo/http"
)

func Init(node *http.RouterNode) {
    node.Middleware(&middleware.HeaderMiddleware{}) // 为当前节点添加上中间件

    node.Get("/", func(request *http.Request, p http.Response) http.ResponseData {
        return http.NewRespDataHtml(404, "blank")
    })

    node.Controller("index", &controller.IndexController{})
}
