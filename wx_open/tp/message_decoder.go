package tp

import (
	"encoding/xml"
	"github.com/bysir-zl/bygo/wx_open/util"
	"github.com/bysir-zl/bygo/wx_open"
)

// 解密消息/事件推送
func DecodeMessage(msgSignature, timeStamp, nonce string, body []byte) (t wx_open.Message, err error) {
	bs, err := util.Decrypt(wx_open.Token, wx_open.AesKey, wx_open.AppId, msgSignature, timeStamp, nonce, body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(bs, &t)
	if err != nil {
		return
	}
	return
}
