package http

import (
	"encoding/json"
	"lib.com/deepzz0/go-com/log"
)

type ResponseData struct {
	Code   int
	Header map[string]string
	Body   string
	Type   string
}

func NewRespDataHtml(code int, body string) ResponseData {
	response := ResponseData{}
	response.Type = "text/html charset=utf-8"

	response.Code = code
	response.Body = body

	return response
}

func NewRespDataJsonString(code int, body string) ResponseData {
	response := ResponseData{}
	response.Type = "application/json charset=utf-8"

	response.Code = code
	response.Body = body

	return response
}

func NewRespDataJson(code int, jsonObj interface{}) ResponseData {
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
