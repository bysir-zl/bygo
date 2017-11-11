package util

import (
    "github.com/bysir-zl/bygo/util/http_util"
    "errors"
)

func Post(url string, body []byte) (rsp []byte, err error) {
    code, rspStr, err := http_util.PostByte(url, body, nil)
    if err != nil {
        return
    }
    if code != 200 {
        err = errors.New("request error,url: " + url+" rsp: "+rspStr)
        return
    }
    rsp = []byte(rspStr)
    return
}
