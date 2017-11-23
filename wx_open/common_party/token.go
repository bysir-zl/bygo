package common_party

import (
	"github.com/bysir-zl/bygo/wx_open/util"
	"fmt"
	"github.com/bysir-zl/bygo/wx_open"
	"encoding/json"
	"sync"
)

type Tokener interface {
	Token() (string, error)
}

type TokenStr string

func (t TokenStr) Token() (string, error) {
	return string(t), nil
}

type Token struct {
	lock sync.Mutex

	appId  string
	secret string
}

type TokenRsp struct {
	wx_open.WxResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (p *Token) Token() (token string, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	key := fmt.Sprintf("at%s", p.appId)
	if x, ok := util.GetData(key); ok {
		token = x.(string)
		return
	}
	tokenRsp, err := GetAccessToken(p.appId, p.secret)
	if err != nil {
		return
	}
	token = tokenRsp.AccessToken

	util.SaveData(key, token, tokenRsp.ExpiresIn)
	return
}

// 获取accessToken, 如果过期则会重新获取
func GetAccessToken(appId string, secret string) (tokenRsp TokenRsp, err error) {
	uri := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	rsp, err := util.Get(fmt.Sprintf(uri, appId, secret))
	if err != nil {
		return
	}

	err = json.Unmarshal(rsp, &tokenRsp)
	if err != nil {
		return
	}
	err = tokenRsp.Error()
	if err != nil {
		return
	}

	return
}

func NewTokener(appId, secret string) Tokener {
	return &Token{
		secret: secret,
		appId:  appId,
	}
}
