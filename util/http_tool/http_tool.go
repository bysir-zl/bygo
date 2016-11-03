package http_tool

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Response string

func (p Response) Json(obj interface{}) {
	json.Unmarshal([]byte(p), obj)
}

func (p Response) String() string {
	return string(p)
}

func Get(url string, params url.Values) (response string, err error) {
	return request(url, "GET", params)
}
func Post(url string, params url.Values) (response string, err error) {
	return request(url, "POST", params)
}

func request(url string, method string, params url.Values) (result string, err error) {
	var response *http.Response

	if method == "GET" {
		up := params.Encode()
		if up != "" {
			if strings.Contains(url, "?") {
				url = url + "&" + up
			} else {
				url = url + "?" + up
			}
		}
		response, err = http.Get(url)
	} else if method == "POST" {
		response, err = http.PostForm(url, params)
	}

	if err != nil {
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
