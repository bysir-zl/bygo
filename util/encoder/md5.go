package encoder

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/bysir-zl/bygo/util"
)

func Md5(in []byte) (out []byte) {
	h := md5.New()
	h.Write(in)
	mbs := h.Sum(nil)
	out = make([]byte, hex.EncodedLen(len(mbs)))
	hex.Encode(out, mbs)

	return
}
func Md5String(in string) (out string) {
	out = util.B2S(Md5(util.S2B(in)))
	return
}
