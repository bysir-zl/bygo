package tp

import (
	"encoding/xml"
	"github.com/bysir-zl/bygo/wx_open"
	"git.coding.net/zzjz/wx_open.git/lib/wx_open"
	"fmt"
)

// 解密消息/事件推送
func DecodeMessage(msgSignature, timeStamp, nonce string, body []byte) (t wx_open.Message, err error) {
	bs, err := util.Decrypt(Token, AesKey, AppId, msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	return
}

// 解密消息/事件推送
func DecodeMessageByte(msgSignature, timeStamp, nonce string, body []byte) (bs []byte, err error) {
	bs, err = util.Decrypt(Token, AesKey, AppId, msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	return
}

// 加密
func EncodeMessageByte(body []byte) (bs []byte, err error) {
	bs, err = util.Encrypt(Token, AesKey, AppId, body)
	return
}

type (
	// 初步解析body, 得到type
	EventReq struct {
		AppId    string `xml:"AppId"`
		InfoType string `xml:"InfoType"`
	}
	EventHandler func(types string, body []byte) (rsp []byte, err error)
)

// 处理授权事件
func HandleEvent(bs []byte, handler EventHandler) (rsp []byte, err error) {
	authMsg := EventReq{}
	err = xml.Unmarshal(bs, &authMsg)
	if err != nil {
		err = fmt.Errorf("unmarshal err:%v; body:%s", err, bs)
		return
	}

	return handler(authMsg.InfoType, bs)
}

type (
	// 初步解析body, 得到type
	MessageReq struct {
		MsgType string `xml:"MsgType"` // 消息类型: event, text ....
	}
	MessageHandler func(msgType string, body []byte) (rsp []byte, err error)
)

// 处理公众号消息
// filterReady: 是否处理全网发布的测试消息
func HandleMessage(bs []byte, appId string, filterReady bool, handler EventHandler) (rsp []byte, err error) {
	msg := MessageReq{}
	err = xml.Unmarshal(bs, &msg)
	if err != nil {
		err = fmt.Errorf("unmarshal err:%v; body:%s", err, bs)
		return
	}

	// 全网发布
	if filterReady {
		var stop bool
		stop, rsp, err = FilterReady(appId, bs)
		if err != nil {
			err = fmt.Errorf("FilterReady err:%v", err)
			return
		}
		if stop {
			return
		}
	}

	return handler(msg.MsgType, bs)
}
