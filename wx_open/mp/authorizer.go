package mp

import (
	"github.com/bysir-zl/bygo/wx_open/util"
	"fmt"
	"errors"
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open"
)

type RefreshUserAccessTokenRsp struct {
	wx_open.WxResponse
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func RefreshUserAccessToken(appid, refreshToken string) (r *RefreshUserAccessTokenRsp, err error) {
	u := "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	rsp, err := util.Get(fmt.Sprintf(u, appid, refreshToken))
	if err != nil {
		err = errors.New("RefreshUserAccessToken err:" + err.Error())
		return
	}
	var rs RefreshUserAccessTokenRsp
	err = json.Unmarshal(rsp, &rs)
	if err != nil {
		return
	}
	err = rs.HasError()
	if err != nil {
		return
	}

	r = &rs
	return
}

type AppAccessTokenRsp struct {
	wx_open.WxResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func GetAppAccessToken(appid, appsecret string) (r *AppAccessTokenRsp, err error) {
	u := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	rsp, err := util.Get(fmt.Sprintf(u, appid, appsecret))
	if err != nil {
		err = errors.New("GetAppAccessToken err:" + err.Error())
		return
	}
	var rs AppAccessTokenRsp
	err = json.Unmarshal(rsp, &rs)
	if err != nil {
		return
	}
	err = rs.HasError()
	if err != nil {
		return
	}

	r = &rs
	return
}
