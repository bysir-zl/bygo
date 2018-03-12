package mp

import (
	"encoding/json"
	"fmt"
	"github.com/bysir-zl/bygo/wx_open/util"
)

type Message struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
}

type MessageText struct {
	Message
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 客服接口-发消息
func SendMessageText(accessToken string, toUser, context string) (err error) {
	m := MessageText{
		Text: struct{ Content string `json:"content"` }{Content: context},
		Message: Message{
			Touser:  toUser,
			Msgtype: "text",
		},
	}
	bs, _ := json.Marshal(m)
	rsp, err := util.Post(URLMessageSend+accessToken, bs)
	if err != nil {
		return
	}

	r := Response{}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return
	}

	if r.ErrCode != 0 {
		err = fmt.Errorf("%+v", r)
		return
	}

	return
}
