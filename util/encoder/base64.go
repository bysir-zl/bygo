package encoder

import (
	"encoding/base64"
	"strings"
	"github.com/bysir-zl/bygo/util"
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
	if src == "" {
		return
	}
	return strings.Trim(util.B2S(Base64Decode(util.S2B(src))), "\x00")
}

func Base64EncodeString(src string) (out string) {
	if src == "" {
		return
	}
	return strings.Trim(util.B2S(Base64Encode(util.S2B(src))), "\x00")
}
