package middleware

import (
    "bygo/http"
)

type HeaderMiddleware struct {

}

func (p HeaderMiddleware) HandlerBefore(s http.SessionContainer) (needStop bool, data http.ResponseData) {
    return false, http.NewRespDataHtml(0, "")
}

func (p HeaderMiddleware) HandlerAfter(s http.SessionContainer) (needStop bool, data http.ResponseData) {

    response := s.Response
    response.AddHeader("Access-Control-Allow-Origin", "*") // 添加上允许跨域
    response.AddHeader("Access-Control-Allow-Headers", "X_TOKEN") // 添加上允许的头,用来身份验证

    return false, http.NewRespDataHtml(0, "")
}
