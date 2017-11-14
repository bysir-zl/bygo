package wxa

import "fmt"

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
	UrlGetCategory   = "https://api.weixin.qq.com/wxa/get_category?access_token="
	UrlGetPage       = "https://api.weixin.qq.com/wxa/get_page?access_token="
	UrlSubmitAudit   = "https://api.weixin.qq.com/wxa/submit_audit?access_token="
	UrlGetLasteAuditstatus   = "https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token="
	UrlRelease   = "https://api.weixin.qq.com/wxa/release?access_token="
	UrlChangeVisitstatus   = "https://api.weixin.qq.com/wxa/change_visitstatus?access_token="
)

type WxResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (p WxResponse) Error() (error) {
	if p.ErrCode == 0 {
		return nil
	}

	return fmt.Errorf("code:%d msg:%s", p.ErrCode, p.ErrMsg)
}
