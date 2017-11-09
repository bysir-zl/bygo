package wx_third_party

import (
	"github.com/bysir-zl/bygo/wx_third_party/util"
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
	c := util.NewCrypt(Token, AesKey, AppId)

	bs, err := c.Decrypt(msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	var t ComponentVerifyTicketReq
	err = xml.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	ticket = t.ComponentVerifyTicket
	err = SaveVerifyTicket(ticket)
	return
}

type SavedVerifyTicket struct {
	VerifyTicket string `json:"verify_ticket"`
	SaveAt       int64  `json:"save_at"`
}

// 获取上一次的ticket, 存储在文件
func GetLastVerifyTicket() (ticket string, ok bool) {
	ks, err := jsonstore.Open("verify_ticket.json")
	if err != nil {
		return
	}

	s := SavedVerifyTicket{}
	err = ks.Get("verify_ticket", &s)
	if err != nil {
		return
	}

	return s.VerifyTicket, true
}

// 存储在文件
func SaveVerifyTicket(ticket string) (err error) {
	ks := new(jsonstore.JSONStore)
	s := SavedVerifyTicket{
		VerifyTicket: ticket,
		SaveAt:       time.Now().Unix(),
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
