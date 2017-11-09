package third_party

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/pkg/errors"
	"github.com/bysir-zl/bygo/wx_open/errs"
	"github.com/bysir-zl/bygo/wx_open"
)

// 第三方平台component_access_token是第三方平台的下文中接口的调用凭据，也叫做令牌（component_access_token）。
// 每个令牌是存在有效期（2小时）的，且令牌的调用不是无限制的，请第三方平台做好令牌的管理，在令牌快过期时（比如1小时50分）再进行刷新。

type ComponentAccessTokenRsp struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int64  `json:"expires_in"`
}

type ComponentAccessTokenReq struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

// 获取ComponentAccessToken
// 会检测过期时间自动刷新哟
func GetComponentAccessToken() (componentAccessToken string, err error) {
	if t, ok := util.GetData("ComponentAccessTokenRsp"); ok {
		return t.(*ComponentAccessTokenRsp).ComponentAccessToken, nil
	}
	ticket, ok := GetLastVerifyTicket()
	if !ok {
		err = errs.ErrNotVerifyTicket
		return
	}

	req := &ComponentAccessTokenReq{
		ComponentVerifyTicket: ticket,
		ComponentAppid:        wx_open.AppId,
		ComponentAppsecret:    wx_open.AppSecret,
	}
	reqData, _ := json.Marshal(req)

	rsp, err := util.Post(URLComponentToken, reqData)
	if err != nil {
		err = errors.Wrap(err, "GetComponentAccessToken")
		return
	}
	var componentAccessTokenRsp ComponentAccessTokenRsp
	err = json.Unmarshal(rsp, &componentAccessTokenRsp)
	if err != nil {
		return
	}
	if componentAccessTokenRsp.ComponentAccessToken == "" {
		err = errors.New(string(rsp))
		return
	}

	util.SaveData("ComponentAccessTokenRsp", &componentAccessTokenRsp, componentAccessTokenRsp.ExpiresIn)

	componentAccessToken = componentAccessTokenRsp.ComponentAccessToken
	return
}
