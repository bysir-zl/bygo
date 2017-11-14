package wxa

import (
	"fmt"
	"github.com/bysir-zl/bygo/wx_open"
	"github.com/bysir-zl/bygo/wx_open/third_party"
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
)

type AuthResponse struct {
	WxResponse
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

// 微信登陆
// code 换取 session_key
// 第三方平台开发者的服务器使用登录凭证 code 以及第三方平台的component_access_token 获取 session_key 和 openid。其中 session_key 是对用户数据进行加密签名的密钥。
// 为了自身应用安全，session_key 不应该在网络上传输。
// appId: 小程序id
func GetSessionKeyByCode(appId string, code string) (r AuthResponse, err error) {
	t, err := third_party.GetComponentAccessToken()
	if err != nil {
		return
	}
	u := fmt.Sprintf(UrlGetSessionKey, appId, code, wx_open.AppId, t)
	rsp, err := util.Post(u, nil)
	if err != nil {
		return
	}

	r = AuthResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}
	return
}
