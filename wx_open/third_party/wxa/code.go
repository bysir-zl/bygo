package wxa

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"fmt"
)

// 代码管理
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1489140610_Uavc4&token=ac82903f643036d2ee5b069f276a6b140a7ab75f&lang=zh_CN

// 为授权的小程序帐号上传小程序代码
// 请注意其中ext_json必须是json的字符串, 而且里面的参数不能乱配置, 比如page参数的路径一定要是小程序模板里有的
func CommitCode(accessToken string, templateId int, extJson string, userVersion string, userDesc string) (err error) {
	req := struct {
		TemplateId  int    `json:"template_id"`
		ExtJson     string `json:"ext_json"`
		UserVersion string `json:"user_version"`
		UserDesc    string `json:"user_desc"`
	}{
		TemplateId:  templateId,
		ExtJson:     extJson,
		UserVersion: userVersion,
		UserDesc:    userDesc,
	}
	reqBs, _ := json.Marshal(req)
	rsp, err := util.Post(UrlCommit+accessToken, reqBs)
	if err != nil {
		return
	}

	r := WxResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}
	return
}

// 获取体验小程序的体验二维码
func GetQrcode(accessToken string) (image []byte, err error) {
	rsp, err := util.Post(UrlGetQrcode+accessToken, nil)
	if err != nil {
		return
	}
	image = rsp

	return
}

type Category struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
	FirstId     int    `json:"first_id"`
	SecondId    int    `json:"second_id"`
}
type CategoryRsp struct {
	WxResponse
	CategoryList []Category `json:"category_list"`
}

// 获取授权小程序帐号的可选类目
func GetCategory(accessToken string) (categoryList []Category, err error) {
	rsp, err := util.Get(UrlGetCategory + accessToken)
	if err != nil {
		return
	}
	r := CategoryRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	categoryList = r.CategoryList
	return
}

type PageRsp struct {
	WxResponse
	PageList []string `json:"page_list"`
}

// 获取小程序的第三方提交代码的页面配置（仅供第三方开发者代小程序调用）
func GetPage(accessToken string) (pages []string, err error) {
	rsp, err := util.Get(UrlGetPage + accessToken)
	if err != nil {
		return
	}
	r := PageRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	pages = r.PageList
	return
}

type AuditItem struct {
	Address     string `json:"address"`
	Tag         string `json:"tag"`
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
	FirstId     int    `json:"first_id"`
	SecondId    int    `json:"second_id"`
	Title       string `json:"title"`
}

type SubmitAuditRsp struct {
	WxResponse
	Auditid int `json:"auditid"`
}

// 将第三方提交的代码包提交审核
func SubmitAudit(accessToken string, auditItems []AuditItem) (auditid int, err error) {
	req, _ := json.Marshal(auditItems)
	rsp, err := util.Post(UrlSubmitAudit+accessToken, req)
	if err != nil {
		return
	}
	r := SubmitAuditRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	auditid = r.Auditid
	return
}

type LatestAuditstatusRsp struct {
	WxResponse
	Auditid string  `json:"auditid"`
	Status  float64 `json:"status"`
	Reason  string  `json:"reason"`
}

// 8. 查询最新一次提交的审核状态
func GetLatestAuditstatus(accessToken string) (r LatestAuditstatusRsp, err error) {
	rsp, err := util.Get(UrlGetLasteAuditstatus + accessToken)
	if err != nil {
		return
	}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}

	return
}

// 9. 发布已通过审核的小程序（仅供第三方代小程序调用）
func Release(accessToken string) (err error) {
	rsp, err := util.Post(UrlRelease+accessToken, []byte("{}"))
	if err != nil {
		return
	}
	r := WxResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}
	return
}

// 10. 修改小程序线上代码的可见状态（仅供第三方代小程序调用）
// action: false:close/true:open
func ChangeVisitstatus(accessToken string, open bool) (err error) {
	action := "close"
	if open {
		action = "open"
	}
	rsp, err := util.Post(UrlChangeVisitstatus+accessToken, []byte(fmt.Sprintf(`{"action":"%s"}`, action)))
	if err != nil {
		return
	}
	r := WxResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.Error()
	if err != nil {
		return
	}
	return
}
