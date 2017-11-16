package common_party

// 生成回复数据

import (
	"encoding/xml"
	"github.com/pkg/errors"
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
func MarshalMessageReply(m interface{}) (bs []byte, err error) {
	_, ok := m.(MessageReply)
	if !ok {
		err = errors.New("bad data, please post a data in MessageReply")
		return
	}
	bs, err = xml.Marshal(m)
	return
}
