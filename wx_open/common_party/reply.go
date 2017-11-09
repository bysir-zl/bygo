package common_party

// 生成回复数据

import (
	"encoding/xml"
	"github.com/bysir-zl/bygo/wx_open"
)

type MessageReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"` // 消息类型
}

type MessageReplyText struct {
	MessageReply
	Content string `xml:"Content"` // 消息内容
}

// 请传入上面定义的结构体
func DecodeMessageReply(m interface{}) (bs []byte, err error) {
	rspBs, _ := xml.Marshal(m)
	bs, err = wx_open.Encrypt(rspBs)
	return
}
