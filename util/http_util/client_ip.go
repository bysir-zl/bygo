package http_util

import (
	"net/http"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/fasthttp-routing"
)

func GetClientIp(r *http.Request) string {
	ip := r.Header.Get("HTTP_X_FORWARDED_FOR")
	if (ip == "") {
		ip = r.RemoteAddr
	}
	return ip
}
func GetClientIpFromFast(p *routing.Context) string {
	ip := util.B2S(p.Request.Header.Peek("HTTP_X_FORWARDED_FOR"))
	if (ip == "") {
		ip = p.RemoteAddr().String()
	}
	return ip
}
