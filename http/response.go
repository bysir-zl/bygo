package http

import (
	"net/http"
)

type Response struct {
	ResponseWrite http.ResponseWriter
	Data          ResponseData

	// todo  result is coding
	Result Result
}

func (p *Response) AddHeader(key string, value string) {
	p.ResponseWrite.Header().Add(key, value)
}

func (p *Response) SetCode(code int) {
	p.ResponseWrite.WriteHeader(code)
}

type Result interface {
	Apply(*http.Request, http.ResponseWriter)
}

