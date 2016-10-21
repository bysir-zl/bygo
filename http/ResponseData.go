package http

import (
    "encoding/json"
)

type ResponseData struct {
    Code      int
    Header    map[string]string;
    Body      string
    Type      string
    Exception Exceptions
}

func NewRespDataHtml(code int, body string) ResponseData {
    response := ResponseData{};

    response.Type = "text/html; charset=utf-8"

    response.Code = code;
    response.Body = body;

    return response;
}

func NewRespDataJson(code int, jsonObj interface{}) ResponseData {
    response := ResponseData{};
    response.Type = "application/json; charset=utf-8"

    response.Code = code;
    bs, err := json.Marshal(jsonObj);
    if err != nil {
        response.Exception = Exceptions{Code:500, Message:"parse to json string fail", Err:err}
    } else {
        response.Body = string(bs);
    }

    return response;
}

func NewRespDataError(code int, err error) ResponseData {

    responseData := ResponseData{Code:200, Exception:Exceptions{Code:code, Err:err, Message:err.Error()}};
    return responseData;
}