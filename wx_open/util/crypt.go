package util

import (
	"github.com/gomydodo/wxencrypter"
)

type Crypt struct {
	e *wxencrypter.Encrypter
}

func NewCrypt(token string, aesKey string, appid string) (c *Crypt, err error) {
	e, err := wxencrypter.NewEncrypter(token, aesKey, appid)
	if err != nil {
		return
	}
	c = &Crypt{
		e: e,
	}
	return
}

// 加密
// replyMsg 公众平台待回复用户的消息，xml格式的字符串,
// bs 包括msg_signature, timestamp, nonce, encrypt的xml格式的字符串.
func (c *Crypt) Encrypt(replyMsg []byte) (bs []byte, err error) {
	// timestamp, nonce 在内部自动生成.
	bs, err = c.e.Encrypt(replyMsg)
	return
}

// 解密
// msgSignature 签名, timestamp 时间戳, nonce 随机, dataMap 要解密的文本
// bs 返回的是解密后的xml格式的字符串
func (c *Crypt) Decrypt(msgSignature, timestamp, nonce string, data []byte) (bs []byte, err error) {
	bs, err = c.e.Decrypt(msgSignature, timestamp, nonce, data)
	return
}

func Decrypt(token, aesKey, appId string, msgSignature, timestamp, nonce string, data []byte) (bs []byte, err error) {
	c, err := NewCrypt(token, aesKey, appId)
	if err != nil {
		return
	}
	return c.Decrypt(msgSignature, timestamp, nonce, data)
}

func Encrypt(token, aesKey, appId string, replyMsg []byte) (bs []byte, err error) {
	c, err := NewCrypt(token, aesKey, appId)
	if err != nil {
		return
	}
	return c.Encrypt(replyMsg)
}
