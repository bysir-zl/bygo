package util

import (
    "github.com/gomydodo/wxencrypter"
)

type Crypt struct {
    token  string
    aesKey string
    appId  string
}

func NewCrypt(token, aesKey, appid string) *Crypt {
    return &Crypt{
        appId:  appid,
        aesKey: aesKey,
        token:  token,
    }
}

// 加密
// replyMsg 公众平台待回复用户的消息，xml格式的字符串,
// bs 包括msg_signature, timestamp, nonce, encrypt的xml格式的字符串.
func (c *Crypt) Encrypt(replyMsg []byte) (bs []byte, err error) {
    e, err := wxencrypter.NewEncrypter(c.token, c.aesKey, c.appId)
    if err != nil {
        return
    }
    // timestamp, nonce 在内部自动生成.
    bs, err = e.Encrypt(replyMsg)

    return
}

// 解密
// msgSignature 签名, timestamp 时间戳, nonce 随机, dataMap 要解密的文本
// bs 返回的是解密后的xml格式的字符串
func (c *Crypt) Decrypt(msgSignature, timestamp, nonce string, data []byte) (bs []byte, err error) {
    e, err := wxencrypter.NewEncrypter(c.token, c.aesKey, c.appId)
    if err != nil {
        return
    }
    bs, err = e.Decrypt(msgSignature, timestamp, nonce, data)
    return
}
