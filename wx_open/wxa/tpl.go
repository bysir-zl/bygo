package wxa

import (
	"fmt"
	"github.com/bysir-zl/bygo/wx_open"
	"encoding/json"
)

type Tpl struct {
	CreateTime  int64  `json:"create_time"`
	UserVersion string `json:"user_version"`
	UserDesc    string `json:"user_desc"`
	DraftId     int    `json:"draft_id"`
}

// 获取草稿箱内的所有临时代码草稿
func GetTplDraft(accessToken string) (tplList []Tpl, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=") + accessToken)
	if err != nil {
		return
	}

	r := struct {
		Errcode      int    `json:"errcode"`
		Errmsg       string `json:"errmsg"`
		TemplateList []Tpl  `json:"template_list"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	tplList = r.TemplateList
	return
}

// 获取代码模版库中的所有小程序代码模版
func GetTpl(accessToken string) (tplList []Tpl, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/gettemplatelist?access_token=") + accessToken)
	if err != nil {
		return
	}

	r := struct {
		Errcode      int    `json:"errcode"`
		Errmsg       string `json:"errmsg"`
		TemplateList []Tpl  `json:"template_list"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	tplList = r.TemplateList
	return
}

// 将草稿箱的草稿选为小程序代码模版
func AddDraftToTpl(accessToken string, draftId int) (err error) {
	req := []byte(fmt.Sprintf(`{"draft_id":%d}`, draftId))
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/addtotemplate?access_token=")+accessToken, req)
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

// 删除指定小程序代码模版
func DelTpl(accessToken string, id int) (err error) {
	req := []byte(fmt.Sprintf(`{"template_id":%d}`, id))
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/deletetemplate?access_token=")+accessToken, req)
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
