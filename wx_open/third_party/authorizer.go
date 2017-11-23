package third_party

import (
	"github.com/bysir-zl/bygo/wx_open/errs"
	"github.com/bysir-zl/bygo/wx_open/util"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"github.com/schollz/jsonstore"
	"time"
	"encoding/xml"
	"sync"
)

const (
	URLComponentToken = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	URLPreAuthCode    = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="
	// 公众号或小程序的接口调用凭据
	URLOtherAuthToken = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="
	// 获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
	URLRefreshOtherAuthToken = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token="
)

var _ = log.Ldate

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

var getComponentAccessTokenLock sync.Mutex

// 获取ComponentAccessToken
// 会检测过期时间自动刷新哟
func GetComponentAccessToken() (componentAccessToken string, err error) {
	getComponentAccessTokenLock.Lock()
	defer getComponentAccessTokenLock.Unlock()

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
		ComponentAppid:        AppId,
		ComponentAppsecret:    AppSecret,
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

// component_verify_ticket
// 出于安全考虑，在第三方平台创建审核通过后，微信服务器每隔10分钟会向第三方的消息接收地址推送一次component_verify_ticket，用于获取第三方平台接口调用凭据
// 接收到后必须直接返回字符串success。

type ComponentVerifyTicketReq struct {
	AppId                 string `xml:"AppId"`
	CreateTime            string `xml:"CreateTime"`
	InfoType              string `xml:"InfoType"`
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
	AuthorizationCode     string `xml:"AuthorizationCode"`
}

// 处理微信VerifyTicket回调
// 成功后会将ticket保存在本地文件
func HandleComponentVerifyTicketReq(msgSignature, timeStamp, nonce string, body []byte) (ticket string, err error) {
	bs, err := util.Decrypt(Token, AesKey, AppId, msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	var t ComponentVerifyTicketReq
	err = xml.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	ticket = t.ComponentVerifyTicket
	if ticket != "" {
		err = SaveVerifyTicket(ticket)
	}
	return
}

// 在内存中缓存一个, 如果服务器重启了这个值为空了, 才重新从文件读取
var stdTicket = ""

type SavedVerifyTicket struct {
	VerifyTicket string `json:"verify_ticket"`
	SaveAt       string `json:"save_at"`
}

// 获取上一次的ticket, 存储在文件
func GetLastVerifyTicket() (ticket string, ok bool) {
	if stdTicket != "" {
		return stdTicket, true
	}

	ks, err := jsonstore.Open("verify_ticket.json")
	if err != nil {
		return
	}

	s := SavedVerifyTicket{}
	err = ks.Get("verify_ticket", &s)
	if err != nil {
		return
	}
	stdTicket = s.VerifyTicket
	if stdTicket == "" {
		return "", false
	}

	return stdTicket, true
}

// 存储在文件
func SaveVerifyTicket(ticket string) (err error) {
	stdTicket = ticket

	ks := new(jsonstore.JSONStore)
	s := SavedVerifyTicket{
		VerifyTicket: ticket,
		SaveAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	err = ks.Set("verify_ticket", s)
	if err != nil {
		return
	}
	err = jsonstore.Save(ks, "verify_ticket.json")
	if err != nil {
		return
	}
	return
}

// 刷新授权公众号或小程序的接口调用凭据（令牌）
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

func RefreshAccessToken(authorizerAppid, refreshToken string) (authorizedInfo *AuthorizedInfoRsp, err error) {
	req := &AuthorizedInfoReq{
		AuthorizerRefreshToken: refreshToken,
		ComponentAppid:         AppId,
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

	authorizedInfo = &r
	return
}

// 该API用于使用授权码换取授权公众号或小程序的授权信息，并换取authorizer_access_token和authorizer_refresh_token。
// 授权码的获取，需要在用户在第三方平台授权页中完成授权流程后，在回调URI中通过URL参数提供给第三方平台方。

type AuthorizerTokenRsp struct {
	AuthorizationInfo struct {
		AuthorizerAppid        string `json:"authorizer_appid"`
		AuthorizerAccessToken  string `json:"authorizer_access_token"`
		ExpiresIn              int64  `json:"expires_in"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
		FuncInfo []struct {
			FuncscopeCategory struct {
				Id int `json:"id"`
			} `json:"funcscope_category"`
		} `json:"func_info"`
	} `json:"authorization_info"`
}

type AuthorizerTokenReq struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

func GetAuthorizerToken(authorizationCode string) (authorizerTokenRsp *AuthorizerTokenRsp, err error) {
	req := &AuthorizerTokenReq{
		ComponentAppid:    AppId,
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

// 该API用于获取预授权码。预授权码用于公众号或小程序授权时的第三方平台方安全验证。

type PreAuthCodeRsp struct {
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int64  `json:"expires_in"`
}
type PreAuthCodeReq struct {
	ComponentAppid string `json:"component_appid"`
}

var getPreAuthCodeLock sync.Mutex
// 获取PreAuthCode
// 会检测过期时间自动刷新哟
func GetPreAuthCode() (preAuthCode string, err error) {
	getPreAuthCodeLock.Lock()
	defer getPreAuthCodeLock.Unlock()

	if t, ok := util.GetData("PreAuthCode"); ok {
		return t.(*PreAuthCodeRsp).PreAuthCode, nil
	}

	componentAccessToken, err := GetComponentAccessToken()

	if err != nil {
		return
	}

	req := &PreAuthCodeReq{
		ComponentAppid: AppId,
	}
	reqData, _ := json.Marshal(req)
	rsp, err := util.Post(URLPreAuthCode+componentAccessToken, reqData)
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
