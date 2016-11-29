package http

import (
	"encoding/json"
	"github.com/deepzz0/go-com/log"
)

type ResponseData struct {
	Code   int
	Header map[string]string
	Body   string
	Type   string
}

func NewResponseHtml(code int, body string) ResponseData {
	response := ResponseData{}
	response.Type = "text/html charset=utf-8"

	response.Code = code
	response.Body = body

	return response
}

func NewResponseJsonString(code int, body string) ResponseData {
	response := ResponseData{}
	response.Type = "application/json charset=utf-8"

	response.Code = code
	response.Body = body

	return response
}

func NewResponseJson(code int, jsonObj interface{}) ResponseData {
	response := ResponseData{}
	response.Type = "application/json charset=utf-8"

	response.Code = code
	bs, err := json.Marshal(jsonObj)
	if err != nil {
		log.Warnf("parse to json string fail : %v", err)
	} else {
		response.Body = string(bs)
	}

	return response
}
