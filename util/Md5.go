package util

import (
    "crypto/md5"
    "encoding/hex"
)

func Md5(in []byte) (out string) {
    h := md5.New();
    h.Write(in);
    out = hex.EncodeToString(h.Sum(nil))

    return
}