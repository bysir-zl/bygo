// 获取授权方的帐号基本信息
// - 公众号
// - 小程序
package tp

import (
	"github.com/bysir-zl/bygo/wx_open"
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
)

type WxAppInfo struct {
	wx_open.WxResponse
	AuthorizerInfo AuthorizerInfo `json:"authorizer_info"`
}

type AuthorizerInfo struct {
	UserName      string `json:"user_name"`      // 小程序的原始ID
	HeadImg       string `json:"head_img"`       // 授权方头像
	NickName      string `json:"nick_name"`      // 授权方昵称
	PrincipalName string `json:"principal_name"` // 小程序的主体名称
	Signature     string `json:"signature"`      // 帐号介绍
	QrcodeUrl     string `json:"qrcode_url"`     // 二维码图片的URL，开发者最好自行也进行保存
}

// 获取小程序账号信息
func GetWxAppInfo(appId string) (wxappInfo WxAppInfo, err error) {
	componentAccessToken, err := GetComponentAccessToken()
	if err != nil {
		return
	}
	req, _ := json.Marshal(map[string]interface{}{
		"component_appid":  AppId,
		"authorizer_appid": appId,
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=")+componentAccessToken, req)
	if err != nil {
		return
	}
	r := WxAppInfo{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	wxappInfo = r
	return
}
