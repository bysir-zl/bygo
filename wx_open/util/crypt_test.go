package util

import (
	"testing"
)

func TestEn(t *testing.T) {
	encodingAesKey := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	token := "pamtest"

	appId := "wxb11529c136998cb6"
	text := `<xml>
	<ToUserName><![CDATA[oia2Tj我是中文jewbmiOUlr6X-1crbLOvLw]]></ToUserName>
	<FromUserName><![CDATA[gh_7f083739789a]]></FromUserName>
	<CreateTime>1407743423</CreateTime>
	<MsgType><![CDATA[video]]></MsgType>
	<Video><
		MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId>
		<Title><![CDATA[testCallBackReplyVideo]]></Title>
		<Description><![CDATA[testCallBackReplyVideo]]></Description>
	</Video>
</xml>`

	c, _ := NewCrypt(token, encodingAesKey, appId)
	x, err := c.Encrypt([]byte(text))
	if err != nil {
		t.Error(t)
	}

	t.Log(string(x))

	//var resXML wxencrypter.EncryptedResponseXML
	//xml.Unmarshal(x, &resXML)
	//encrypt := resXML.Encrypt
	//msgSignature := resXML.MsgSignature
	//format := "<xml><ToUserName><![CDATA[toUser]]></ToUserName><Encrypt><![CDATA[%s]]></Encrypt></xml>"
	//fromXML := fmt.Sprintf(format, encrypt)
	//timestamp := "1409304348"
	//nonce := "xxxxxx"
	//
	//bs,err:=c.Decrypt(msgSignature, timestamp, nonce, []byte(fromXML))
	//if err != nil {
	//    t.Error(err)
	//}
	//
	//t.Log(string(bs))

}
