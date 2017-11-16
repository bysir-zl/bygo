package third_party

import (
	"github.com/bysir-zl/bygo/wx_open/common_party"
	"time"
	"strings"
	"github.com/bysir-zl/bygo/wx_open"
	"github.com/bysir-zl/bygo/wx_open/util"
)

// 1. 模拟粉丝触发专用测试公众号的事件，并推送事件消息到专用测试公众号，第三方平台方开发者需要提取推送XML信息中的event值，并在5秒内立即返回按照下述要求组装的文本消息给粉丝。
func responseMock1(req wx_open.Message) (bs []byte, err error) {
	rsp := common_party.MessageReplyText{
		MessageReply: common_party.MessageReply{
			MsgType:      "text",
			ToUserName:   req.FromUserName,
			FromUserName: req.ToUserName,
			CreateTime:   time.Now().Unix(),
		},

		Content: req.Event + "from_callback",
	}

	bs, err = common_party.MarshalMessageReply(rsp)
	if err != nil {
		return
	}
	bs, err = util.Encrypt(Token, AesKey, AppId, bs)
	return
}

// 2. 模拟粉丝发送文本消息给专用测试公众号，第三方平台方需根据文本消息的内容进行相应的响应：
func responseMock2(req wx_open.Message) (bs []byte, err error) {
	rsp := common_party.MessageReplyText{
		MessageReply: common_party.MessageReply{
			MsgType:      "text",
			ToUserName:   req.FromUserName,
			FromUserName: req.ToUserName,
			CreateTime:   time.Now().Unix(),
		},
		Content: "TESTCOMPONENT_MSG_TYPE_TEXT_callback",
	}

	bs, err = common_party.MarshalMessageReply(rsp)
	if err != nil {
		return
	}
	bs, err = util.Encrypt(Token, AesKey, AppId, bs)
	return
}

// 3. 模拟粉丝发送文本消息给专用测试公众号，第三方平台方需在5秒内返回空串表明暂时不回复，然后再立即使用客服消息接口发送消息回复粉丝
func responseMock3(req wx_open.Message) (err error) {
	if !strings.Contains(req.Content, "QUERY_AUTH_CODE") {
		return
	}

	authCode := strings.Replace(req.Content, "QUERY_AUTH_CODE:", "", -1)
	t, err := GetAuthorizerToken(authCode)
	if err != nil {
		return
	}

	at := t.AuthorizationInfo.AuthorizerAccessToken
	err = common_party.SendMessageText(at, req.FromUserName, authCode+"_from_api")
	if err != nil {
		return
	}
	return
}

// 拦截测试消息并相应
func FilterReady(appId string, req wx_open.Message) (stop bool, response []byte, err error) {
	if appId != "wx570bc396a51b8ff8" && appId != "wxd101a85aa106f53e" {
		return
	}

	stop = true
	switch req.MsgType {
	case "event":
		response, err = responseMock1(req)
		return
	case "text":
		if req.Content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
			response, err = responseMock2(req)
			return
		}
		if strings.Contains(req.Content, "QUERY_AUTH_CODE") {
			err = responseMock3(req)
			return
		}
	}

	return
}
