package controller

import (
    "fmt"
    "strings"
    "github.com/bysir-zl/bygo/http"
)

type IndexController struct{}

func (p IndexController) Index(r *http.Request, s http.Response) http.ResponseData {

    return http.NewRespDataHtml(200, "welcome to use bygo!" + "<br>" +
        "Url: " + r.Router.Url + "<br>" +
        "Handler: " + r.Router.Handler + "<br>" +
        "RouterParams : " + strings.Join(r.Router.Params, ",") + "<br>" +
        "Input : " + fmt.Sprint(r.Input.All()) + "<br>" +
        "Header : " + fmt.Sprint(r.Header) + "<br>" +
        "")
}
