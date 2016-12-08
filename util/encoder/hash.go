package encoder

import (
	"crypto/sha256"
	"encoding/hex"
	"crypto"
	"github.com/bysir-zl/bygo/util"
)

func Sha256(in string) string {
	hash := sha256.New()
	hash.Write([]byte(in))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}

func Hash(in []byte, hash crypto.Hash) (out []byte) {
	h := hash.New()
	h.Write(in)
	md := h.Sum(nil)
	out = make([]byte, hex.EncodedLen(len(md)))
	hex.Encode(out, md)
	return
}

func HashString(in string, hash crypto.Hash) (out string) {
	h := hash.New()
	h.Write(util.S2B(in))
	md := h.Sum(nil)
	out = hex.EncodeToString(md)
	return
}