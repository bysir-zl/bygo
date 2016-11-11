package encoder

import (
	"encoding/base64"
	"strings"
)

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)

	return buf
}

func Base64Decode(src []byte) (out []byte) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(buf, src)

	return buf
}

func Base64DecodeString(src string) (out string) {
	return strings.Trim(string(Base64Decode([]byte(src))), "\x00")
}
