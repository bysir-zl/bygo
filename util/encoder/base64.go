package encoder

import (
	"encoding/base64"
	"github.com/bysir-zl/bygo/util"
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

// 在go中, 如果解码的字符串缺少最后的=号, 将不能解码, 所以先填补缺少的=
func Base64DecodeString(src string) (out string) {
	if src == "" {
		return
	}
	a := len(src) % 3
	if a != 0 {
		src = src + strings.Repeat("=", 3 - a)
	}
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}
	out = util.B2S(bs)
	return
}

func Base64EncodeString(src string) (out string) {
	if src == "" {
		return
	}
	out = base64.StdEncoding.EncodeToString(util.S2B(src))
	return
}
