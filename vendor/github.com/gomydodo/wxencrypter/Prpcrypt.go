package wxencrypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"time"
)

type Prpcrypt struct {
	key     []byte
	Encoder PKCS7Encoder
}

func NewPrpcrypt(key string) (p *Prpcrypt, err error) {
	b, err := base64.StdEncoding.DecodeString(key + "=")
	if err != nil {
		err = DecodeBase64Error
		return
	}
	p = &Prpcrypt{
		key:     b,
		Encoder: pkcs7Encoder{},
	}
	return
}

func (p *Prpcrypt) Encrypt(appId string, src []byte) (ret string, err error) {
	b, err := aes.NewCipher(p.key)
	if err != nil {
		err = EncryptAESError
		return
	}

	buf := &bytes.Buffer{}
	random := p.random()

	_, err = buf.Write(random)
	if err != nil {
		err = EncryptAESError
		return
	}

	err = binary.Write(buf, binary.BigEndian, int32(len(src)))
	if err != nil {
		err = EncryptAESError
		return
	}

	_, err = buf.Write(src)

	if err != nil {
		err = EncryptAESError
		return
	}
	_, err = buf.WriteString(appId)

	if err != nil {
		err = EncryptAESError
		return
	}
	content := buf.Bytes()

	content = p.Encoder.Encode(content)

	c := cipher.NewCBCEncrypter(b, p.key[:16])
	c.CryptBlocks(content, content)
	ret = base64.StdEncoding.EncodeToString(content)
	return
}

func (p *Prpcrypt) Decrypt(appId string, src string) (ret []byte, err error) {
	b, err := aes.NewCipher(p.key)
	if err != nil {
		err = DecryptAESError
		return
	}

	content, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		err = DecodeBase64Error
		return
	}

	c := cipher.NewCBCDecrypter(b, p.key[:16])
	c.CryptBlocks(content, content)

	content = p.Encoder.Decode(content)[16:]

	xmlLen := binary.BigEndian.Uint32(content[:4])

	ret = content[4 : xmlLen+4]

	fromAppId := string(content[xmlLen+4:])
	if appId != fromAppId {
		err = ValidateAppidError
		return
	}

	return
}

func (p *Prpcrypt) random() (b []byte) {
	src := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")
	n := len(src)
	buf := &bytes.Buffer{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 16; i += 1 {
		index := r.Intn(n)
		buf.WriteByte(src[index])
	}
	b = buf.Bytes()
	return
}
