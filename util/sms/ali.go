// 阿里的短信服务

package sms

import (
	"github.com/bysir-zl/bygo/util"
	"strconv"
	"time"
	"encoding/json"
	"net/url"
	"github.com/bysir-zl/bygo/util/encoder"
	"crypto"
	"github.com/bysir-zl/bygo/util/http_util"
	"errors"
	"strings"
)

const (
	HostAliSms = "http://dysmsapi.aliyuncs.com"
	ApiAliSms  = ""
)

type Ali struct {
	apiKey    string
	apiSecret string
}

func NewAli(apiKey, apiSecret string) *Ali {
	return &Ali{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// 发送短信
// phones 支持以逗号,分割的多个手机号码
func (a *Ali) Send(tplCode string, signName string, phones string, data map[string]string) (error) {
	bs, _ := json.Marshal(data)
	l, _ := time.LoadLocation("GMT")
	timestamp := time.Now().In(l).Format("2006-01-02T15:04:05Z")
	params := map[string]string{
		// 系统参数
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   strconv.Itoa(util.Rand(0, 9999999)),
		"AccessKeyId":      a.apiKey,
		"SignatureVersion": "1.0",
		"Timestamp":        timestamp,
		"Format":           "JSON",
		// 业务API参数
		"Action":        "SendSms",
		"Version":       "2017-05-25",
		"RegionId":      "cn-hangzhou",
		"PhoneNumbers":  phones,
		"SignName":      signName,
		"TemplateParam": string(bs),
		"TemplateCode":  tplCode,
		"OutId":         "123",
	}
	kv := util.ParseOrderKV(params)

	sign, _ := a.Sign(kv)
	kv.Add("Signature", sign)

	code, rsp, err := http_util.Get(HostAliSms+ApiAliSms, kv, nil)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New("rsp status is't 200, rsp:" + rsp)
	}
	if !strings.Contains(rsp, "OK") {
		return errors.New("rsp :" + rsp)
	}

	return nil
}

// POP签名
func (a *Ali) Sign(kv util.OrderKV) (sign string, err error) {
	kv.Sort()
	signStr := "GET" + "&" + url.QueryEscape("/") + "&" + url.QueryEscape(kv.EncodeString())
	signBs := encoder.Hmac([]byte(signStr), []byte(a.apiSecret+"&"), crypto.SHA1)
	signBs = encoder.Base64Encode(signBs)
	sign = string(signBs)
	return
}
