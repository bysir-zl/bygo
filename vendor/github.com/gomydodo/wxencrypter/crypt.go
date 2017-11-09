package wxencrypter

import (
	"fmt"
	"time"
)

type Encrypter struct {
	prpcrypter     *Prpcrypt
	token          string
	encodingAesKey string
	appId          string
}

func NewEncrypter(token, encodingAesKey, appId string) (e *Encrypter, err error) {
	if len(encodingAesKey) != 43 {
		err = IllegalAesKey
		return
	}

	p, err := NewPrpcrypt(encodingAesKey)
	if err != nil {
		return
	}

	e = &Encrypter{
		prpcrypter:     p,
		token:          token,
		appId:          appId,
		encodingAesKey: encodingAesKey,
	}
	return
}

func (e *Encrypter) Encrypt(replyMsg []byte) (b []byte, err error) {
	encrypt, err := e.prpcrypter.Encrypt(e.appId, replyMsg)
	if err != nil {
		return
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := string(e.prpcrypter.random())
	signature := Sha1(e.token, timestamp, nonce, encrypt)

	b, err = GenerateResponseXML(encrypt, signature, timestamp, nonce)
	return
}

func (e *Encrypter) Decrypt(msgSignature, timestamp, nonce string, data []byte) (b []byte, err error) {
	reqXML, err := ParseRequestXML(data)
	if err != nil {
		return
	}

	signature := Sha1(e.token, timestamp, nonce, reqXML.Encrypt)
	if signature != msgSignature {
		err = ValidateSignatureError
		return
	}
	b, err = e.prpcrypter.Decrypt(e.appId, reqXML.Encrypt)
	return
}
