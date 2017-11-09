package wxencrypter

import (
	"bytes"
)

type PKCS7Encoder interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
}

const (
	blockSize = 32
)

type pkcs7Encoder struct {
}

func (p pkcs7Encoder) Encode(src []byte) (dist []byte) {
	byteLen := len(src)
	pad := blockSize - (byteLen % blockSize)
	if pad == 0 {
		pad = blockSize
	}

	b := bytes.Repeat([]byte{byte(pad)}, pad)
	dist = append(src, b...)
	return
}

func (p pkcs7Encoder) Decode(src []byte) (dist []byte) {
	byteLen := len(src)
	pad := int(src[byteLen-1])
	if pad < 1 || pad > 32 {
		dist = src
	} else {
		dist = src[:byteLen-pad]
	}
	return
}
