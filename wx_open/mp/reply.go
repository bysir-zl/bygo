// 生成回复数据

package mp

import (
	"encoding/xml"
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

func (x *MessageReplyText) Byte() []byte {
	bs, _ := xml.Marshal(x)
	return bs
}
