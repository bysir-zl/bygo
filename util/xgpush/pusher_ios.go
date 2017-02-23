package xgpush

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/bysir-zl/bygo/util"
	"strconv"
)

type PusherIos struct {
	PusherBase
}

func NewPusherIos(accessId, accessKey, secretKey string) Pusher {
	return &PusherIos{
		PusherBase{
			access_id:  accessId,
			secret_key: secretKey,
			access_key: accessKey,
		},
	}
}

// push消息给单个设备
// typ: 消息类型：1：通知 2：透传消息。iOS 平台请填 0
func (p *PusherIos) PushSingleDriver(token string, typ int, sendTime string, message Message) (err error) {
	ps := util.OrderKV{}
	ps.Add("device_token", token)
	ps.Add("message_type", "0")
	ps.Add("message", message.StringIos())
	ps.Add("expire_time", strconv.Itoa(3*24*3600)) // 3 天
	if sendTime != "" {
		ps.Add("send_time", sendTime) // year-mon-day hour:min:sec 若小于服务器当前时间，则会立即推送
	}
	e := "2"
	if beego.BConfig.RunMode == "dev" {
		e = "1"
	}
	ps.Add("environment", e) // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台丌填或填 0

	err = p.requestServer("push/single_device", ps, nil)

	return
}

// push消息给单个用户
// typ: 消息类型：1：通知 2：透传消息。iOS 平台请填 0
func (p *PusherIos) PushSingleAccount(account string, typ int, sendTime string, message Message) (err error) {
	ps := util.OrderKV{}
	ps.Add("account", account)
	ps.Add("message_type", "0")
	ps.Add("message", message.StringIos())
	if sendTime != "" {
		ps.Add("send_time", sendTime) // year-mon-day hour:min:sec 若小于服务器当前时间，则会立即推送
	}
	ps.Add("expire_time", strconv.Itoa(3*24*3600)) // 3 天
	e := "2"
	if beego.BConfig.RunMode == "dev" {
		e = "1"
	}
	ps.Add("environment", e) // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台丌填或填 0

	err = p.requestServer("push/single_account", ps, nil)

	return
}

// push消息给全部
// typ: 消息类型：1：通知 2：透传消息。iOS 平台请填 0
func (p *PusherIos) PushAll(typ int, sendTime string, message Message) (err error) {
	ps := util.OrderKV{}
	ps.Add("message_type", "0")
	ps.Add("message", message.StringIos())
	if sendTime != "" {
		ps.Add("send_time", sendTime) // year-mon-day hour:min:sec 若小于服务器当前时间，则会立即推送
	}
	ps.Add("expire_time", strconv.Itoa(3*24*3600)) // 3 天
	e := "2"
	if beego.BConfig.RunMode == "dev" {
		e = "1"
	}
	ps.Add("environment", e) // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台丌填或填 0

	err = p.requestServer("push/all_device", ps, nil)

	return
}

// push消息给tags
// tagsOpType 1:AND 2:OR
func (p *PusherIos) PushByTags(tags []string, tagsOpType string, typ int, sendTime string, message Message) (err error) {
	ps := util.OrderKV{}
	ps.Add("message_type", "0")
	ps.Add("message", message.StringIos())
	ts, _ := json.Marshal(&tags)
	ps.Add("tags_list", string(ts))
	if len(tags) == 1 {
		tagsOpType = "OR"
	}
	ps.Add("tags_op", tagsOpType)
	if sendTime != "" {
		ps.Add("send_time", sendTime) // year-mon-day hour:min:sec 若小于服务器当前时间，则会立即推送
	}
	ps.Add("expire_time", strconv.Itoa(3*24*3600)) // 3 天
	e := "2"
	if beego.BConfig.RunMode == "dev" {
		e = "1"
	}
	ps.Add("environment", e) // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台丌填或填 0

	err = p.requestServer("push/tags_device", ps, nil)

	return
}
