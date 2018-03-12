package tp_wxa

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
)

type DomainRsp struct {
	wx_open.WxResponse
	Requestdomain   []string `json:"requestdomain"`
	Wsrequestdomain []string `json:"wsrequestdomain"`
	Uploaddomain    []string `json:"uploaddomain"`
	Downloaddomain  []string `json:"downloaddomain"`
}

// 修改服务器地址
// 授权给第三方的小程序，其服务器域名只可以为第三方的服务器，当小程序通过第三方发布代码上线后，
// 小程序原先自己配置的服务器域名将被删除， 只保留第三方平台的域名，所以第三方平台在代替小程序发布代码之前，
// 需要调用接口为小程序添加第三方自身的域名。
// action: add添加, delete删除, set覆盖, get获取。当参数是get时不需要填四个域名字段。
func ModifyDomain(accessToken string, action string, requestDomain, wsrequestdomain, uploaddomain, downloaddomain []string) (r DomainRsp, err error) {
	req := map[string]interface{}{
		"action":          action,
		"requestdomain":   requestDomain,
		"wsrequestdomain": wsrequestdomain,
		"uploaddomain":    uploaddomain,
		"downloaddomain":  downloaddomain,
	}

	reqBs, _ := json.Marshal(req)
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/modify_domain?access_token=")+accessToken, reqBs)
	if err != nil {
		return
	}

	r = DomainRsp{}
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
