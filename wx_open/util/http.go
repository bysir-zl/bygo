package util

import (
	"errors"
	"crypto/tls"
	"net/http"
	"time"
	"bytes"
	"io/ioutil"
)

func Post(url string, body []byte) (rsp []byte, err error) {
	code, rsp, err := request("POST", url, body)
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New("request error,url: " + url + " rsp: " + string(rsp))
		return
	}
	return
}

func Get(url string) (rsp []byte, err error) {
	code, rsp, err := request("GET", url, nil)
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New("request error,url: " + url + " rsp: " + string(rsp))
		return
	}

	return
}

func request(method string, url string, postBody []byte) (code int, rsp []byte, err error) {
	// 忽略https证书验证
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest("GET", url, nil)
	} else if method == "POST" {
		req, _ = http.NewRequest("POST", url, bytes.NewReader(postBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	response, err := client.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()
	code = response.StatusCode
	rsp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}
