package huyi_sms

import (
	"errors"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/http_util"
)

type Config struct {
	ApiId  string
	ApiKey string
}

type Sms struct {
	c *Config
}

const ApiHost = "http://106.ihuyi.com/webservice/sms.php"

func NewSms(c *Config) *Sms {
	return &Sms{
		c: c,
	}
}

func (p *Sms) Send(phone, content string) (err error) {
	ps := util.OrderKV{}
	ps.Add("method", "Submit")
	ps.Add("account", p.c.ApiId)
	ps.Add("password", p.c.ApiKey)
	ps.Add("mobile", phone)
	ps.Add("content", content)
	ps.Add("format", "json")
	_, rsp, err := http_util.Get(ApiHost, ps, nil)
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
