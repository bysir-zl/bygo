package tp_wxa

import (
	"fmt"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open/util"
	"encoding/json"
	"log"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
)

var _ = log.Ldate

// 消息模板管理

type MessageTpl struct {
	TemplateId string `json:"template_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Example    string `json:"example"`
}

// 获取帐号下已存在的模板列表
func GetMsgTpl(accessToken string, offset int, count int) (tplList []MessageTpl, err error) {
	req := []byte(fmt.Sprintf(`{"offset":%d,"count":%d}`, offset, count))

	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/wxopen/template/list?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		Errcode int          `json:"errcode"`
		Errmsg  string       `json:"errmsg"`
		List    []MessageTpl `json:"list"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	tplList = r.List
	return
}

// 删除帐号下的某个模板
func DelMsgTpl(accessToken string, id string) (err error) {
	req := []byte(fmt.Sprintf(`{"template_id":"%s"}`, id))

	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/wxopen/template/del?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	return
}

/*ps: 关于模板添加:
微信写得接口文档太不详细, 没写清楚逻辑, 以至于有接口 但不知道接口作用. 在这重新整理说明:
发送模板通知必须得到模板id, 而这个模板不保存在第三方平台中, 而是通过接口向公众平台添加消息模板, 得到模板id后再发送.
通过第三方管理公众平台的接口微信有写, 比较容易理解, 这里详细说下怎么添加模板
微信写了三个接口来实现添加模板:
1.获取小程序模板库标题列表
2.获取模板库某个模板标题下关键词库
3.组合模板并添加至帐号下的个人模板库
无任何逻辑解释, 只有猜着理解.

去研究公众平台的[模板消息]才理解:
1. 模板只能选择在微信模板库已有的模板, 而怎么得到已有模板库? 就是接口1
2. 模板里的KV也只能使用模板库里已有的, 比如商城支付成功通知 里面只能有订单号:xxx等信息, 而不能有管道疏通:158xxxxxxxx等信息, 怎么得到KV列表? 就是接口2
3. 得到模板标题, 和模板KV, 就能组装一个完整的模板并添加了, 这就是 接口3所说的组合模板
*/

type WxTplTitle struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// 获取小程序模板库标题列表
// accessToken: 公众平台的token, 而不是第三方的token
func GetTplMsgLibTitleList(accessToken string, offset int, count int) (tplList []WxTplTitle, total int64, err error) {
	req := []byte(fmt.Sprintf(`{"offset":%d,"count":%d}`, offset, count))

	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/wxopen/template/library/list?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		Errcode    int          `json:"errcode"`
		Errmsg     string       `json:"errmsg"`
		List       []WxTplTitle `json:"list"`
		TotalCount int64        `json:"total_count"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	tplList = r.List
	total = r.TotalCount
	return
}

type Keyword struct {
	KeywordId int    `json:"keyword_id"`
	Name      string `json:"name"`
	Example   string `json:"example"`
}

// 获取模板库某个模板标题下关键词库
func GetMsgTplLibKV(accessToken string, tplTitleId string) (kvs []Keyword, err error) {
	req := []byte(fmt.Sprintf(`{"id":"%s"}`, tplTitleId))

	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/wxopen/template/library/get?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		Errcode     int       `json:"errcode"`
		Errmsg      string    `json:"errmsg"`
		KeywordList []Keyword `json:"keyword_list"`
		Title       string    `json:"title"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	kvs = r.KeywordList
	return
}

// 组合模板并添加至帐号下的个人模板库
func AddWxTplToSelf(accessToken string, tplTitleId string, keywordIdList []int) (templateId string, err error) {
	reqB := map[string]interface{}{
		"id":              tplTitleId,
		"keyword_id_list": keywordIdList,
	}
	req, _ := json.Marshal(reqB)

	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/wxopen/template/add?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		wx_open.WxResponse
		TemplateId string `json:"template_id"`
	}{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.HasError()
	if err != nil {
		return
	}

	templateId = r.TemplateId
	return
}

type MsgTplArg struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// touser 是 接收者（用户）的 openid
// template_id 是 所需下发的模板消息的id
// page 否 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
// form_id 是 表单提交场景下，为 submit 事件带上的 formId；支付场景下，为本次支付的 prepay_id
// data 是 模板内容，不填则下发空模板
// color 否 模板内容字体的颜色，不填默认黑色
// emphasis_keyword 否 模板需要放大的关键词，不填则默认无放大
func SendTpl(accessToken string, toUserOpenId string, tplId string, formId string, page string, color string, data map[string]MsgTplArg, emphasisKeyword string) (err error) {
	reqB := map[string]interface{}{
		"touser":      toUserOpenId,
		"template_id": tplId,
		"form_id":     formId,
	}
	if page != "" {
		reqB["page"] = page
	}
	if data != nil {
		reqB["data"] = data
	}
	if emphasisKeyword != "" {
		reqB["emphasis_keyword"] = emphasisKeyword
	}
	if color != "" {
		reqB["color"] = color
	}

	req, _ := json.Marshal(reqB)
	rsp, err := util.Post(("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=")+accessToken, req)
	if err != nil {
		return
	}

	r := wx_open.WxResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.HasError()
	if err != nil {
		return
	}

	return
}
