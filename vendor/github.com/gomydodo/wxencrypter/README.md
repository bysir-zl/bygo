#微信公众号开发消息体签名及加密解密方案的GO语言实现

	import(
		"github.com/gomydodo/wxencrypter"
		"encoding/xml"
		"fmt"
	)

	func main() {
		encodingAesKey := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
		token := "pamtest"
		timestamp := "1409304348"
		nonce := "xxxxxx"
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

		e, err := wxencrypter.NewEncrypter(token, encodingAesKey, appId)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := e.Encrypt([]byte(text), timestamp, nonce)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))

		var resXML wxencrypter.EncryptedResponseXML
		xml.Unmarshal(b, &resXML)
		encrypt := resXML.Encrypt
		msgSignature := resXML.MsgSignature
		format := "<xml><ToUserName><![CDATA[toUser]]></ToUserName><Encrypt><![CDATA[%s]]></Encrypt></xml>"
		fromXML := fmt.Sprintf(format, encrypt)
		b, err = e.Decrypt(msgSignature, timestamp, nonce, []byte(fromXML))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
	}


###消息体加密

*微信正常消息为明文，格式如下：*

	msg = 
	<xml>
		<ToUserName><![CDATA[toUser]]></ToUserName>
		<FromUserName><![CDATA[fromUser]]></FromUserName> 
		<CreateTime>1348831860</CreateTime>
		<MsgType><![CDATA[text]]></MsgType>
		<Content><![CDATA[this is a test]]></Content>
		<MsgId>1234567890123456</MsgId>
	</xml>

*公众号选择了安全模式后，消息格式如下*， 调用`e.Decrypt(msgSignature, timestamp, nonce, new_msg)`解密

	new_msg=
	<xml> 
		<ToUserName><![CDATA[toUser]]</ToUserName>
	    <Encrypt><![CDATA[msg_encrypt]]</Encrypt>
	</xml>


微信服务器请求时在url上增加参数：msg_signature


| 参数 | 说明  |
| :-------- | ----------: |
|  Token  | 公众平台上，开发者设置的Token |
|  timestamp  | URL上原有参数，时间戳 |
|  nonce | URL上原有参数，随机数 |
|  msg_encrypt | 前文描述密文消息体 |


### 回复消息体的签名与加密

*明文格式：*

	msg=
	<xml>
		 <ToUserName><![CDATA[toUser]]></ToUserName>
		 <FromUserName><![CDATA[fromUser]]></FromUserName>
		 <CreateTime>12345678</CreateTime>
		 <MsgType><![CDATA[text]]></MsgType>
		 <Content><![CDATA[你好]]></Content>
	</xml>

调用`e.Encrypt(msg, timestamp, nonce)`加密

*加密后消息格式：*

	new_msg=
	<xml>
		<Encrypt><![CDATA[msg_encrypt]]></Encrypt>
		<MsgSignature><![CDATA[msg_signature]]></MsgSignature>
		<TimeStamp>timestamp</TimeStamp>
		<Nonce><![CDATA[nonce]]></Nonce>
	</xml> 