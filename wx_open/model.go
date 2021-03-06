package wx_open

import "fmt"

// 更多资料请看微信公众平台文档: https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140454
// 事件消息结构体
type Message struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`   // 消息类型: event, text ....
	Event        string `xml:"Event"`     // 事件类型
	EventKey     string `xml:"EventKey"`  // 事件KEY值
	Ticket       string `xml:"Ticket"`    // 二维码的ticket，可用来换取二维码图片
	Latitude     string `xml:"Latitude"`  // 地理位置纬度
	Longitude    string `xml:"Longitude"` // 地理位置经度
	Precision    string `xml:"Precision"` // 地理位置精度

	Content string `xml:"Content"` // 消息内容
	MsgId   string `xml:"MsgId"`
	Reason  string `xml:"Reason"`
}

// 已知的接口都是这个错误格式, 可以统一处理
type WxResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WxResponseErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (p WxResponseErr) Error() string {
	return fmt.Sprintf("code:%d msg:%s", p.ErrCode, p.ErrMsg)
}

func (p WxResponse) HasError() (error) {
	if p.ErrCode == 0 {
		return nil
	}

	return WxResponseErr{
		ErrCode: p.ErrCode,
		ErrMsg:  p.ErrMsg,
	}
}
