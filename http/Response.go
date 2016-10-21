package http

import (
    "net/http"
)
type Response struct {
    ResponseWrite http.ResponseWriter
    ResponseData ResponseData
}

func (p *Response)AddHeader(key string,value string)  {
    p.ResponseWrite.Header().Add(key,value)
}

func (p *Response)SetCode(code int)  {
    p.ResponseWrite.WriteHeader(code)
}

