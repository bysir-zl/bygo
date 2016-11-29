package http

import (
	"net/http"
)

type Response struct {
	ResponseWrite http.ResponseWriter
	Data          ResponseData


	Result Result
}

func (p *Response) AddHeader(key string, value string) {
	p.ResponseWrite.Header().Add(key, value)
}

func (p *Response) SetCode(code int) {
	p.ResponseWrite.WriteHeader(code)
}

// todo  result is coding , now is not use
type Result interface {
	Apply(*http.Request, http.ResponseWriter)
}

