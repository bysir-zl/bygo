package xgpush

import (
	"encoding/json"
	"fmt"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/encoder"
	"github.com/bysir-zl/bygo/util/http_util"
	"strconv"
	"strings"
	"time"
)

type PusherBase struct {
	access_id string
	access_key string
	secret_key string
}


// 给token设置tag
// tags [][tag,token]
func (p *PusherBase) SetTags4Token(tags [][2]string) (err error) {
	ps := util.OrderKV{}

	t, _ := json.Marshal(&tags)
	ps.Add("tag_token_list", string(t))

	err = p.requestServer("tags/batch_set", ps, nil)
	return
}

// 给批量删除tag
// tags [][tag,token]
func (p *PusherBase) DelTags4Token(tags [][2]string) (err error) {
	ps := util.OrderKV{}

	t, _ := json.Marshal(&tags)
	ps.Add("tag_token_list", string(t))

	err = p.requestServer("tags/batch_del", ps, nil)
	return
}

// 获取token的tags
func (p *PusherBase) GetTokenTags(token string) (err error, tags []string) {
	ps := util.OrderKV{}
	ps.Add("device_token", token)

	type Ts struct {
		Tags []string `json:"tags,omitempty"`
	}
	ts := Ts{}
	err = p.requestServer("tags/query_token_tags", ps, &ts)
	tags = ts.Tags

	return
}

func (p *PusherBase) requestServer(api string, ps util.OrderKV, rsp interface{}) (err error) {
	method := "POST"
	url := fmt.Sprintf("http://%s/v2/%s", host, api)

	ps.Add("access_id", p.access_id)
	ps.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	ps.Sort()

	// sign
	keys := ps.Keys()
	values := ps.Values()
	u := strings.Replace(url, "http://", "", -1)
	signString := method + u
	for i, key := range keys {
		value := values[i]
		signString += key + "=" + value
	}

	sign := encoder.Md5String(signString + p.secret_key)

	ps.Add("sign", sign)

	log.Info("request", url+"?"+ps.EncodeStringWithoutEscape())
	// request
	_, rs, err := http_util.Post(url, ps, nil)
	if err != nil {
		return
	}
	log.Info("response", rs)

	bj, err := bjson.New([]byte(rs))
	if err != nil {
		return
	}
	if code := bj.Pos("ret_code").Int(); code != 0 {
		msg := bj.Pos("err_msg").String()
		err = fmt.Errorf("code is %d, msg is %s", code, msg)
		return
	}

	bjr := bj.Pos("result")
	if rsp != nil && !bjr.IsNil() {
		err = bjr.Object(rsp)
	}
	return
}


