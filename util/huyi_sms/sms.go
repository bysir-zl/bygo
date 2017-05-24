package huyi_sms

import (
	"errors"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/http_util"
	"strings"
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
	_, rsp, err := http_util.Get(ApiHost, ps, nil)
	if err != nil {
		return
	}

	if !strings.Contains(rsp, "<code>2</code>") {
		errInfo := rsp
		msg := strings.Split(rsp, "</msg>")[0]
		if strings.Contains(msg, "<msg>") {
			errInfo = strings.Split(msg, "<msg>")[1]
		}
		err = errors.New(errInfo)
		return
	}

	return
}
