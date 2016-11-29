package http_util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"bytes"
	"github.com/deepzz0/go-com/log"
	"crypto/tls"
	"github.com/bysir-zl/bygo/util"
)

type Response string

func (p Response) Json(obj interface{}) {
	json.Unmarshal([]byte(p), obj)
}

func (p Response) String() string {
	return string(p)
}

func Get(url string, params url.Values, header map[string]string) (response string, err error) {
	bs, err := request(url, "GET", params, header)
	response = util.B2S(bs)
	return
}
func Post(url string, params url.Values, header map[string]string) (response string, err error) {
	bs, err := request(url, "POST", params, header)
	response = util.B2S(bs)
	return
}

func request(url string, method string, params url.Values, header map[string]string) (result []byte, err error) {
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
	//if response.StatusCode != 200 {
	//	err = errors.New("request stausCode is " + strconv.Itoa(response.StatusCode))
	//	return
	//}
	body, _ := ioutil.ReadAll(response.Body)
	result = body
	return

}

func BuildQuery(key []string, value []string) string {
	s := ""
	for i, k := range key {
		s = s + "&" + k + "=" + value[i]
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

func BuildQueryWithOutEmptyValue(key []string, value []string) string {
	s := ""
	for i, k := range key {
		if value[i] == "" {
			continue
		}
		s = s + "&" + k + "=" + value[i]
	}
	if s != "" {
		s = s[1:]
	}
	return s
}


func QueryString2Map(que string) (set map[string]string) {
	set = map[string]string{}
	if !strings.Contains(que, "&") {
		return
	}
	for _, kv := range strings.Split(que, "&") {
		kAv := strings.Split(kv, "&")
		if len(kAv) == 2 {
			k, err := url.QueryUnescape(kAv[0])
			v, err2 := url.QueryUnescape(kAv[1])
			if err == nil && err2 == nil {
				set[k] = v
			}
		}
	}
	return
}