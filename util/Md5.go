package util

import (
    "crypto/md5"
    "encoding/hex"
)

func Md5(in string) (out string) {
    h := md5.New();
    h.Write([]byte(in));
    out = hex.EncodeToString(h.Sum(nil))

    return
}