package third_party

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/pkg/errors"
	"github.com/bysir-zl/bygo/wx_open"
)

// 该API用于使用授权码换取授权公众号或小程序的授权信息，并换取authorizer_access_token和authorizer_refresh_token。
// 授权码的获取，需要在用户在第三方平台授权页中完成授权流程后，在回调URI中通过URL参数提供给第三方平台方。

type AuthorizerTokenRsp struct {
	AuthorizationInfo struct {
		AuthorizerAppid        string `json:"authorizer_appid"`
		AuthorizerAccessToken  string `json:"authorizer_access_token"`
		ExpiresIn              int64  `json:"expires_in"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
		FuncInfo struct {
			FuncscopeCategory struct {
				Id string `json:"id"`
			} `json:"funcscope_category"`
		}
	} `json:"authorization_info"`
}

type AuthorizerTokenReq struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

func GetAuthorizerToken(authorizationCode string) (authorizerTokenRsp *AuthorizerTokenRsp, err error) {
	req := &AuthorizerTokenReq{
		ComponentAppid:    wx_open.AppId,
		AuthorizationCode: authorizationCode,
	}
	reqData, _ := json.Marshal(req)

	componentAccessToken, err := GetComponentAccessToken()
	if err != nil {
		return
	}
	rsp, err := util.Post(URLOtherAuthToken+componentAccessToken, reqData)
	if err != nil {
		err = errors.Wrap(err, "GetAuthorizerToken")
		return
	}

	authorizerTokenRsp = &AuthorizerTokenRsp{}
	err = json.Unmarshal(rsp, authorizerTokenRsp)
	if err != nil {
		return
	}
	if authorizerTokenRsp.AuthorizationInfo.AuthorizerAppid == "" {
		err = errors.New(string(rsp))
		return
	}

	key := "AccessToken:" + authorizerTokenRsp.AuthorizationInfo.AuthorizerAppid
	util.SaveData(key, authorizerTokenRsp.AuthorizationInfo.AuthorizerAccessToken, authorizerTokenRsp.AuthorizationInfo.ExpiresIn)

	return
}
