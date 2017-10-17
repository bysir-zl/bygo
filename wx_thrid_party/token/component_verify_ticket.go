package token

import (
    "github.com/bysir-zl/bygo/wx_thrid_party/config"
    "github.com/bysir-zl/bygo/wx_thrid_party/util"
    "encoding/xml"
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

// 根据微信回调解析出VerifyTicket
func DecodeComponentVerifyTicket(msg_signature, timeStamp, nonce string, body []byte) (ticket string, err error) {
    c := util.NewCrypt(config.Token, config.AesKey, config.AppId)

    bs, err := c.Decrypt(msg_signature, timeStamp, nonce, body)
    if err != nil {
        return
    }

    var t ComponentVerifyTicketReq
    err = xml.Unmarshal(bs, &t)
    ticket = t.ComponentVerifyTicket
    return
}

// 获取上一次的ticket, 存储在文件
func GetLastVerifyTicket() (ticket string, ok bool) {
    return "xx", true
}
