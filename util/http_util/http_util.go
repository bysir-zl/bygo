package http_util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"bytes"
	"lib.com/deepzz0/go-com/log"
	"crypto/tls"
)

type Response string

func (p Response) Json(obj interface{}) {
	json.Unmarshal([]byte(p), obj)
}

func (p Response) String() string {
	return string(p)
}

func Get(url string, params url.Values, header map[string]string) (response string, err error) {
	return request(url, "GET", params, header)
}
func Post(url string, params url.Values, header map[string]string) (response string, err error) {
	return request(url, "POST", params, header)
}

func request(url string, method string, params url.Values, header map[string]string) (result string, err error) {
	var response *http.Response

	// 忽略https证书验证
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}
	var req *http.Request
	if method == "GET" {
		up := params.Encode()
		if up != "" {
			if strings.Contains(url, "?") {
				url = url + "&" + up
			} else {
				url = url + "?" + up
			}
		}
		req, _ = http.NewRequest("GET", url, nil)
	} else if method == "POST" {
		req, _ = http.NewRequest("POST", url, bytes.NewReader([]byte(params.Encode())))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if header != nil&&len(header) != 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	response, err = client.Do(req)
	if err != nil {
		log.Warn("http request error : ", err)
		return
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		result = string(body)
		return
	}
	err = errors.New("request stausCode is " + strconv.Itoa(response.StatusCode))
	return

}
