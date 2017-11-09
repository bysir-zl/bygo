package wxencrypter

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"sort"
)

func Sha1(token, timestamp, nonce, msg string) (ret string) {
	sl := []string{token, timestamp, nonce, msg}
	sort.Strings(sl)

	h := sha1.New()
	for _, s := range sl {
		io.WriteString(h, s)
	}
	encode := h.Sum(nil)

	ret = hex.EncodeToString(encode)
	return
}
