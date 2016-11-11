package http_util

import "net/http"

func GetClientIp(r *http.Request) string {
	ip := r.Header.Get("HTTP_X_FORWARDED_FOR")
	if (ip == "") {
		ip = r.RemoteAddr
	}
	return ip
}
