package third_party

import (
	"fmt"
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/bysir-zl/bygo/wx_open"
	"errors"
)

const (
	UrlCommit        = "https://api.weixin.qq.com/wxa/commit?access_token="
	UrlModifyDomain  = "https://api.weixin.qq.com/wxa/modify_domain?access_token="
	UrlGetQrcode     = "https://api.weixin.qq.com/wxa/get_qrcode?access_token="
	UrlGetTplDraft   = "https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token="
	UrlGetTpl        = "https://api.weixin.qq.com/wxa/gettemplatelist?access_token="
	UrlAddToTpl      = "https://api.weixin.qq.com/wxa/addtotemplate?access_token="
	UrlDelTpl        = "https://api.weixin.qq.com/wxa/deletetemplate?access_token="
	UrlGetSessionKey = "https://api.weixin.qq.com/sns/component/jscode2session?appid=%s&js_code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s"
	UrlGetMsgTpl     = "https://api.weixin.qq.com/cgi-bin/wxopen/template/list?access_token="
	UrlDelMsgTpl     = "https://api.weixin.qq.com/cgi-bin/wxopen/template/del?access_token="
	UrlGetWxTplList  = "https://api.weixin.qq.com/cgi-bin/wxopen/template/library/list?access_token="
	UrlGetWxTplKV    = "https://api.weixin.qq.com/cgi-bin/wxopen/template/library/get?access_token="
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

	r := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.ErrCode, r.ErrMsg)
		return
	}
	return
}

// 修改服务器地址
// 授权给第三方的小程序，其服务器域名只可以为第三方的服务器，当小程序通过第三方发布代码上线后，
// 小程序原先自己配置的服务器域名将被删除， 只保留第三方平台的域名，所以第三方平台在代替小程序发布代码之前，
// 需要调用接口为小程序添加第三方自身的域名。
func ModifyDomain(accessToken string, action string, requestDomain, wsrequestdomain, uploaddomain, downloaddomain []string) (err error) {
	req := map[string]interface{}{
		"action":          action,
		"requestdomain":   requestDomain,
		"wsrequestdomain": wsrequestdomain,
		"uploaddomain":    uploaddomain,
		"downloaddomain":  downloaddomain,
	}

	reqBs, _ := json.Marshal(req)
	rsp, err := util.Post(UrlModifyDomain+accessToken, reqBs)
	if err != nil {
		return
	}

	r := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.ErrCode, r.ErrMsg)
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

type Tpl struct {
	CreateTime  int64  `json:"create_time"`
	UserVersion string `json:"user_version"`
	UserDesc    string `json:"user_desc"`
	DraftId     int    `json:"draft_id"`
}

// 获取草稿箱内的所有临时代码草稿
func GetTplDraft(accessToken string) (tplList []Tpl, err error) {
	rsp, err := util.Get(UrlGetTplDraft + accessToken)
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
	rsp, err := util.Get(UrlGetTpl + accessToken)
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
	rsp, err := util.Post(UrlAddToTpl+accessToken, req)
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
	rsp, err := util.Post(UrlDelTpl+accessToken, req)
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

type AuthResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

// 微信登陆
// code 换取 session_key
// 第三方平台开发者的服务器使用登录凭证 code 以及第三方平台的component_access_token 获取 session_key 和 openid。其中 session_key 是对用户数据进行加密签名的密钥。
// 为了自身应用安全，session_key 不应该在网络上传输。
// appId: 小程序id
func GetSessionKeyByCode(appId string, code string) (r AuthResponse, err error) {
	t, err := GetComponentAccessToken()
	if err != nil {
		return
	}
	u := fmt.Sprintf(UrlGetSessionKey, appId, code, wx_open.AppId, t)
	rsp, err := util.Post(u, nil)
	if err != nil {
		return
	}

	r = AuthResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Openid == "" {
		err = errors.New(string(rsp))
		return
	}
	return
}

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

	rsp, err := util.Post(UrlGetMsgTpl+accessToken, req)
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

	rsp, err := util.Post(UrlDelMsgTpl+accessToken, req)
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
是不是不知所云? 无任何逻辑解释, 只有猜着理解...

去研究公众平台的[模板消息]才理解到:
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
func GetWxTplTitleList(accessToken string, offset int, count int) (tplList []WxTplTitle, total int64, err error) {
	req := []byte(fmt.Sprintf(`{"offset":%d,"count":%d}`, offset, count))

	rsp, err := util.Post(UrlGetWxTplList+accessToken, req)
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
func GetWxTplKV(accessToken string, tplTitleId string) (kvs []Keyword, err error) {
	req := []byte(fmt.Sprintf(`{"id":"%s"}`, tplTitleId))

	rsp, err := util.Post(UrlGetWxTplKV+accessToken, req)
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

	rsp, err := util.Post(UrlGetWxTplKV+accessToken, req)
	if err != nil {
		return
	}

	r := struct {
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
		TemplateId string `json:"template_id"`
	}{}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = fmt.Errorf("code: %d msg:%s", r.Errcode, r.Errmsg)
		return
	}

	templateId = r.TemplateId
	return
}
