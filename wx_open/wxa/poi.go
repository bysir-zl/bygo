package wxa

import (
	"git.coding.net/zzjz/wx_open.git/lib/wx_open/util"
	"encoding/json"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
)

// 地点审核消息
type PoiAudioMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"` // 消息类型: event, text ....
	Event        string `xml:"Event"`
	AuditId      int64  `xml:"audit_id"` // 审核单id
	Status       int    `xml:"status"`   // 审核状态（3：审核通过，2：审核失败）
	Reason       string `xml:"reason"`   // 如果status为3或者4，会返回审核失败的原因
	PoiId        int64  `xml:"poi_id"`
}

type AddNearByPoiRsp struct {
	wx_open.WxResponse
	Data struct {
		AuditId           int64  `json:"audit_id"`
		PoiId             int64  `json:"poi_id"`
		RelatedCredential string `json:"related_credential"`
	} `json:"data"`
}

// https://mp.weixin.qq.com/debug/wxadoc/dev/api/nearby.html#添加地点
func AddNearByPoi(token string, relatedName, relatedCredential, relatedAddress, relatedProofMaterial string) (*AddNearByPoiRsp, error) {
	req, _ := json.Marshal(map[string]interface{}{
		"related_name":           relatedName,
		"related_credential":     relatedCredential,
		"related_address":        relatedAddress,
		"related_proof_material": relatedProofMaterial,
	})
	rsp, err := util.Post(("https://api.weixin.qq.com/wxa/addnearbypoi?access_token=")+token, req)
	if err != nil {
		return nil, err
	}

	r := AddNearByPoiRsp{}
	err = json.Unmarshal(rsp, r)
	if err != nil {
		return nil, err
	}
	if e := r.HasError(); e != nil {
		return nil, e
	}
	return &r, nil
}
