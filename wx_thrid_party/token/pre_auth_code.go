package token

import (
    "encoding/json"
    "github.com/bysir-zl/bygo/wx_thrid_party/config"
    "github.com/bysir-zl/bygo/wx_thrid_party/util"
    "github.com/pkg/errors"
    "github.com/bysir-zl/bygo/wx_thrid_party/errs"
)

// 该API用于获取预授权码。预授权码用于公众号或小程序授权时的第三方平台方安全验证。

type PreAuthCodeRsp struct {
    PreAuthCode string `json:"pre_auth_code"`
    ExpiresIn   int64  `json:"expires_in"`
}
type PreAuthCodeReq struct {
    ComponentAppid string `json:"component_appid"`
}

// 获取PreAuthCode
// 会检测过期时间自动刷新哟
func GetPreAuthCode() (preAuthCode string, err error) {
    if t, ok := util.GetData("PreAuthCode"); ok {
        return t.(*PreAuthCodeRsp).PreAuthCode, nil
    }

    req := &PreAuthCodeReq{
        ComponentAppid: config.AppId,
    }
    reqData, _ := json.Marshal(req)

    componentAccessToken, err := GetComponentAccessToken()
    if err != nil {
        return
    }
    rsp, err := util.Post(URL_PreAuthcode+componentAccessToken, reqData)
    if err != nil {
        err = errors.Wrap(err, "GetPreAuthCode")
        return
    }
    var preAuthCodeRsp PreAuthCodeRsp
    err = json.Unmarshal(rsp, &preAuthCodeRsp)
    if err != nil {
        return
    }

    util.SaveData("PreAuthCode", &preAuthCodeRsp, preAuthCodeRsp.ExpiresIn)

    preAuthCode = preAuthCodeRsp.PreAuthCode
    return
}
