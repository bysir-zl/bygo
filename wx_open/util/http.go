package util

import (
	"github.com/bysir-zl/bygo/util/http_util"
	"errors"
	"github.com/bysir-zl/bygo/util"
)

func Post(url string, body []byte) (rsp []byte, err error) {
	code, rspStr, err := http_util.PostByte(url, body, nil)
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New("request error,url: " + url + " rsp: " + rspStr)
		return
	}
	rsp = []byte(rspStr)
	return
}

func Get(url string) (rsp []byte, err error) {
	code, rspStr, err := http_util.Get(url, util.OrderKV{}, nil)
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New("request error,url: " + url + " rsp: " + rspStr)
		return
	}
	rsp = []byte(rspStr)
	return
}
