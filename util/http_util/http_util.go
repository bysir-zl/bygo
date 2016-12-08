package http_util

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"bytes"
	"github.com/deepzz0/go-com/log"
	"crypto/tls"
	"github.com/bysir-zl/bygo/util"
)

func Get(url string, params util.OrderKV, header map[string]string) (response string, err error) {
	up := params.EncodeString()
	if up != "" {
		if strings.Contains(url, "?") {
			url = url + "&" + up
		} else {
			url = url + "?" + up
		}
	}
	bs, err := request(url, "GET", nil, header, nil)
	response = util.B2S(bs)
	return
}

func Post(url string, params util.OrderKV, header map[string]string) (response string, err error) {
	bs, err := request(url, "POST", params.Encode(), header, nil)
	response = util.B2S(bs)
	return
}

func PostByte(url string, post []byte, header map[string]string) (response string, err error) {
	bs, err := request(url, "POST", post, header, nil)
	response = util.B2S(bs)
	return
}

func PostWithCookie(url string, params util.OrderKV, cookie map[string]string) (response string, err error) {
	bs, err := request(url, "POST", params.Encode(), nil, cookie)
	response = util.B2S(bs)
	return
}

func GetWithCookie(url string, params util.OrderKV, cookie map[string]string) (response string, err error) {
	bs, err := request(url, "GET", params.Encode(), nil, cookie)
	response = util.B2S(bs)
	return
}

func request(url string, method string, post []byte, header map[string]string, cookie map[string]string) (result []byte, err error) {
	var response *http.Response

	// 忽略https证书验证
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}
	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest("GET", url, nil)
	} else if method == "POST" {
		req, _ = http.NewRequest("POST", url, bytes.NewReader(post))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if header != nil&&len(header) != 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	if cookie != nil&&len(cookie) != 0 {
		for key, value := range cookie {
			req.AddCookie(&http.Cookie{Name:key, Value:value})
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
	var bf bytes.Buffer
	for i, k := range key {
		if bf.Len() == 0 {
			bf.WriteByte('&')
		}
		bf.WriteString(k + "=" + url.QueryEscape(value[i]))
	}
	return bf.String()
}

func BuildQueryWithOutEmptyValue(key []string, value []string) string {
	var bf bytes.Buffer
	for i, k := range key {
		if v := value[i]; v != "" {
			if bf.Len() == 0 {
				bf.WriteByte('&')
			}
			bf.WriteString(k + "=" + url.QueryEscape(v))
		}
	}
	return bf.String()
}

func QueryString2Map(que string) (set map[string]string) {
	set = map[string]string{}
	if !strings.Contains(que, "&") {
		return
	}
	for _, kv := range strings.Split(que, "&") {
		kAv := strings.Split(kv, "=")
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

// like php - rawurlencode
// rawurlencode and urlencode is different form the ' ' will encode to '%20', is not '+'
func RawUrlEncode(origin string) string {
	x := url.QueryEscape(origin)
	x = strings.Replace(x, "+", "%20", -1)
	return x
}