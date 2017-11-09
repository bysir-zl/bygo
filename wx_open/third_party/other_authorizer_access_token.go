package third_party

import (
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/bysir-zl/bygo/wx_open"
)

// 获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
// 该API用于在授权方令牌（authorizer_access_token）失效时，可用刷新令牌（authorizer_refresh_token）获取新的令牌

type AuthorizedInfoReq struct {
	ComponentAppid         string `json:"component_appid"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}
type AuthorizedInfoRsp struct {
	AccessToken  string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"authorizer_refresh_token"`
}

func GetOrRefreshAccessToken(componentAppid, authorizerAppid, refreshToken string) (authorizedInfo *AuthorizedInfoRsp, isNew bool, err error) {
	key := "AccessToken:" + authorizerAppid
	if t, ok := util.GetData(key); ok {
		acT := t.(string)
		authorizedInfo = &AuthorizedInfoRsp{
			AccessToken: acT,
		}
		return authorizedInfo, false, nil
	}

	req := &AuthorizedInfoReq{
		AuthorizerRefreshToken: refreshToken,
		ComponentAppid:         wx_open.AppId,
		AuthorizerAppid:        authorizerAppid,
	}
	reqData, _ := json.Marshal(req)

	componentAccessToken, err := GetComponentAccessToken()
	if err != nil {
		return
	}

	rsp, err := util.Post(URLRefreshOtherAuthToken+componentAccessToken, reqData)
	if err != nil {
		err = errors.Wrap(err, "RefreshAccessToken")
		return
	}
	var r AuthorizedInfoRsp
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}

	util.SaveData(key, r.AccessToken, r.ExpiresIn)

	authorizedInfo = &r
	isNew = true
	return
}
