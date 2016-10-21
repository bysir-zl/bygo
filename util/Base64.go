package util

import "encoding/base64"

func Base64Encode(in string) string {
    return base64.StdEncoding.EncodeToString([]byte(in))
}

func Base64Decode(in string) (out string) {
    bs,err := base64.StdEncoding.DecodeString(in);
    if err != nil {
        return
    }
    return string(bs)
}
