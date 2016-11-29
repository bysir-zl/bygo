package encoder

import (
	"crypto/hmac"
	"crypto/sha1"
)

func HmacSha1(origin, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(origin)
	return mac.Sum(nil)
}