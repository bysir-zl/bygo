package encoder

import (
	"crypto/hmac"
	"crypto"
)

func Hmac(origin, key []byte, hash crypto.Hash) []byte {
	mac := hmac.New(hash.New, key)
	mac.Write(origin)
	return mac.Sum(nil)
}