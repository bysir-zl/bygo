package util

import (
    "crypto/sha256"
    "encoding/hex"
)

func Sha256(in string) string {

    hash := sha256.New()
    hash.Write([]byte(in))
    md := hash.Sum(nil)
    mdStr := hex.EncodeToString(md)
    return mdStr
}
