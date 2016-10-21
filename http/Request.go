package http

import (
    "net/http"
    "bygo/http/input"
)

// 保存当前请求的路由信息
type router struct {
    Params  []string

    Url     string
    Hash    string
    Handler string
}

type Request struct {
    Request *http.Request
    Router  router
    Input   input.Input
    Method  string
    Header  map[string]string
}

func (p *Request)Init() {
    inpt := input.Input{};
    inpt.Request = p.Request;

    //method
    p.Method = p.Request.Method;

    //header
    p.Header = make(map[string]string);
    for key, value := range p.Request.Header {
        p.Header[key] = value[0]
    }

    p.Input = inpt;
}
