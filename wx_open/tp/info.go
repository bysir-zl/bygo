// 获取授权方的帐号基本信息
// - 公众号
// - 小程序
package tp

import (
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open/util"
	"encoding/json"
	"errors"
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
	VerifyTypeInfo struct {
		Id int `json:"id"`
	} `json:"verify_type_info"`                  // 授权方认证类型，-1代表未认证，0代表微信认证
	MiniProgramInfo struct {
		Network struct {
			RequestDomain   []string `json:"RequestDomain"`
			WsRequestDomain []string `json:"WsRequestDomain"`
			UploadDomain    []string `json:"UploadDomain"`
			DownloadDomain  []string `json:"DownloadDomain"`
		} `json:"network"`
	} `json:"MiniProgramInfo"`
	BusinessInfo struct {
		OpenPay   int `json:"open_pay"`   // 是否开通微信支付功能
		OpenShake int `json:"open_shake"` // 是否开通微信摇一摇功能
		OpenScan  int `json:"open_scan"`  // 是否开通微信扫商品功能
		OpenCard  int `json:"open_card"`  // 是否开通微信卡券功能
		OpenStore int `json:"open_store"` // 是否开通微信门店功能
	} `json:"business_info"`
}

// 获取小程序账号信息
func GetWxAppInfo(appId string) (wxappInfo AuthorizerInfo, err error) {
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
		err = errors.New("json Unmarshal err:" + err.Error())
		return
	}
	err = r.HasError()
	if err != nil {
		return
	}

	wxappInfo = r.AuthorizerInfo
	return
}
