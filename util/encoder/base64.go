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
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}
	return util.B2S(bs)
}

func Base64EncodeString(src string) (out string) {
	if src == "" {
		return
	}
	out = base64.StdEncoding.EncodeToString(util.S2B(src))
	return
}
