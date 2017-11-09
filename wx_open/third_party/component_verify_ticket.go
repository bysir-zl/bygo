package third_party

import (
	"github.com/bysir-zl/bygo/wx_open"
	"encoding/xml"
	"time"
	"github.com/schollz/jsonstore"
)

// component_verify_ticket
// 出于安全考虑，在第三方平台创建审核通过后，微信服务器每隔10分钟会向第三方的消息接收地址推送一次component_verify_ticket，用于获取第三方平台接口调用凭据
// 接收到后必须直接返回字符串success。
type ComponentVerifyTicketReq struct {
	AppId                 string `xml:"AppId"`
	CreateTime            string `xml:"CreateTime"`
	InfoType              string `xml:"InfoType"`
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
	AuthorizationCode     string `xml:"AuthorizationCode"`
}

// 处理微信VerifyTicket回调
// 成功后会将ticket保存在本地文件
func HandleComponentVerifyTicketReq(msgSignature, timeStamp, nonce string, body []byte) (ticket string, err error) {
	bs, err := wx_open.Decrypt(msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	var t ComponentVerifyTicketReq
	err = xml.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	ticket = t.ComponentVerifyTicket
	if ticket != "" {
		err = SaveVerifyTicket(ticket)
	}
	return
}

// 在内存中缓存一个, 如果服务器重启了这个值为空了, 才重新从文件读取
var stdTicket = ""

type SavedVerifyTicket struct {
	VerifyTicket string `json:"verify_ticket"`
	SaveAt       string `json:"save_at"`
}

// 获取上一次的ticket, 存储在文件
func GetLastVerifyTicket() (ticket string, ok bool) {
	if stdTicket != "" {
		return ticket, true
	}

	ks, err := jsonstore.Open("verify_ticket.json")
	if err != nil {
		return
	}

	s := SavedVerifyTicket{}
	err = ks.Get("verify_ticket", &s)
	if err != nil {
		return
	}
	stdTicket = s.VerifyTicket
	if stdTicket == "" {
		return "", false
	}

	return s.VerifyTicket, true
}

// 存储在文件
func SaveVerifyTicket(ticket string) (err error) {
	stdTicket = ticket

	ks := new(jsonstore.JSONStore)
	s := SavedVerifyTicket{
		VerifyTicket: ticket,
		SaveAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	err = ks.Set("verify_ticket", s)
	if err != nil {
		return
	}
	err = jsonstore.Save(ks, "verify_ticket.json")
	if err != nil {
		return
	}
	return
}
