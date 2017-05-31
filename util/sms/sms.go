package sms

import (
	"errors"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/http_util"
	"strings"
)

type HuyiSms struct {
	ApiId  string
	ApiKey string
}

const ApiHost_Huyi = "http://106.ihuyi.com/webservice/sms.php"

func NewHuyiSms(ApiId, ApiKey string) *HuyiSms {
	return &HuyiSms{
		ApiId: ApiId, ApiKey: ApiKey,
	}
}

func (p *HuyiSms) Send(phone, content string) (err error) {
	ps := util.OrderKV{}
	ps.Add("method", "Submit")
	ps.Add("account", p.ApiId)
	ps.Add("password", p.ApiKey)
	ps.Add("mobile", phone)
	ps.Add("content", content)
	ps.Add("format", "json")
	_, rsp, err := http_util.Get(ApiHost_Huyi, ps, nil)
	if err != nil {
		return
	}

	bj, err := bjson.New([]byte(rsp))
	if err != nil {
		err = errors.New(rsp)
		return
	}
	if bj.Pos("code").Int() != 2 {
		err = errors.New(bj.Pos("msg").String())
		return
	}

	return
}

/*----------------------*/

const Apidaiyi = "http://api.daiyicloud.com/asmx/smsservice.aspx"

type Daiyi struct {
	Name string
	Pwd  string
}

func NewDaiyi(name, pwd string) *Daiyi {
	return &Daiyi{
		Name: name, Pwd: pwd,
	}
}

func (p *Daiyi) Send(mobile string, content, sign string) (err error) {
	ps := util.OrderKV{}
	ps.Add("name", p.Name)
	ps.Add("pwd", p.Pwd)
	ps.Add("content", content)
	ps.Add("mobile", mobile)
	ps.Add("sign", sign)
	ps.Add("type", "pt")

	_, rsp, err := http_util.Post(Apidaiyi, ps, nil)
	if err != nil {
		return
	}

	// code,sendid,invalidcount,successcount,blackcount,msg
	ds := strings.Split(rsp, ",")
	if len(ds) != 6 {
		err = errors.New(rsp)
		return
	}
	if ds[0] != "0" {
		err = errors.New(rsp)
		return
	}

	return
}
