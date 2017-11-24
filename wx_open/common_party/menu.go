package common_party

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/bysir-zl/bygo/wx_open"
)

type Button struct {
	Type       string   `json:"type,omitempty"`       // 非必须; 菜单的响应动作类型
	Name       string   `json:"name,omitempty"`       // 必须;  菜单标题
	Key        string   `json:"key,omitempty"`        // 非必须; 菜单KEY值, 用于消息接口推送
	URL        string   `json:"url,omitempty"`        // 非必须; 网页链接, 用户点击菜单可打开链接
	MediaId    string   `json:"media_id,omitempty"`   // 非必须; 调用新增永久素材接口返回的合法media_id
	AppId      string   `json:"appid,omitempty"`      // 非必须; 跳转到小程序的appid
	PagePath   string   `json:"pagepath,omitempty"`   // 非必须; 跳转到小程序的path
	SubButtons []Button `json:"sub_button,omitempty"` // 非必须; 二级菜单数组
}

// 个性化菜单规则
// see https://mp.weixin.qq.com/wiki?action=doc&id=mp1455782296&t=0.8648400919429371#1
type MatchRule struct {
	GroupId            string `json:"group_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
	TagId              string `json:"tag_id,omitempty"`
}

type Menu struct {
	Buttons   []Button   `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	MenuId    int64      `json:"menuid,omitempty"`
}

// 创建
// see https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141013
func CreateMenu(accessToken Tokener, menu *Menu) (error) {
	t, err := accessToken.Token()
	if err != nil {
		return err
	}
	uri := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + t
	req, _ := json.Marshal(menu)

	rsp, err := util.Post(uri, req)
	if err != nil {
		return err
	}
	r := wx_open.WxResponse{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}

	return nil
}
