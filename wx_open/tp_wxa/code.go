package tp_wxa

import (
	"encoding/json"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open/util"
	"fmt"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
)

// 代码管理
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1489140610_Uavc4&token=ac82903f643036d2ee5b069f276a6b140a7ab75f&lang=zh_CN

// 为授权的小程序帐号上传小程序代码
// 请注意其中ext_json必须是json的字符串, 而且里面的参数不能乱配置, 比如page参数的路径一定要是小程序模板里有的
func CommitCode(accessToken string, templateId int64, extJson string, userVersion string, userDesc string) (err error) {
	req := struct {
		TemplateId  int64  `json:"template_id"`
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
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/commit?access_token=")+accessToken, reqBs)
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

type PageRsp struct {
	wx_open.WxResponse
	PageList []string `json:"page_list"`
}

// 获取小程序的第三方提交代码的页面配置（仅供第三方开发者代小程序调用）
func GetPage(accessToken string) (pages []string, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/get_page?access_token=") + accessToken)
	if err != nil {
		return
	}
	r := PageRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.HasError()
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
	wx_open.WxResponse
	Auditid int64 `json:"auditid"`
}

// 将第三方提交的代码包提交审核
func SubmitAudit(accessToken string, auditItems []AuditItem) (auditid int64, err error) {
	req, _ := json.Marshal(map[string]interface{}{
		"item_list": auditItems,
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/submit_audit?access_token=")+accessToken, req)
	if err != nil {
		return
	}
	r := SubmitAuditRsp{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	err = r.HasError()
	if err != nil {
		return
	}

	auditid = r.Auditid
	return
}

type LatestAuditstatusRsp struct {
	wx_open.WxResponse
	Auditid string  `json:"auditid"`
	Status  float64 `json:"status"`
	Reason  string  `json:"reason"`
}

// 8. 查询最新一次提交的审核状态
func GetLatestAuditstatus(accessToken string) (r LatestAuditstatusRsp, err error) {
	rsp, err := util.Get(("https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token=") + accessToken)
	if err != nil {
		return
	}

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

// 9. 发布已通过审核的小程序（仅供第三方代小程序调用）
func Release(accessToken string) (err error) {
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/release?access_token=")+accessToken, []byte("{}"))
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

// 10. 修改小程序线上代码的可见状态（仅供第三方代小程序调用）
// action: false:close/true:open
func ChangeVisitstatus(accessToken string, open bool) (err error) {
	action := "close"
	if open {
		action = "open"
	}
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/change_visitstatus?access_token=")+accessToken, []byte(fmt.Sprintf(`{"action":"%s"}`, action)))
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
