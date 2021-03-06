package auth

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/util/encoder"
	"strings"
	"time"
)

type JWTData struct {
	Iss string
	Iat int64 // iat(issued at): 在什么时候签发的
	Exp int64
	Sub string
	Aud string
	Typ string
}

func (j JWTData) IsEmpty() bool {
	return j.Iat == 0
}

// 加密数据
// iss: 该JWT的签发者 ,exp(expires): 什么时候过期，这里是一个Unix时间戳 ,typ:用户类型(组),如admin/user ,sub: 该JWT所面向的用户 ,aud: 接收该JWT的一方
func JWTEncode(secret,iss string, exp int64, typ string, sub string, aud string) (out string) {
	mapper := JWTData{}
	mapper.Iss = iss
	mapper.Iat = time.Now().Unix()
	mapper.Exp = exp
	mapper.Typ = typ
	mapper.Sub = sub
	mapper.Aud = aud

	bs, e := json.Marshal(mapper)
	if e != nil {
		return
	}

	payload := encoder.Base64EncodeString(string(bs))
	if payload == "" {
		return
	}

	header := encoder.Base64EncodeString("{\"typ\":\"JWT\",\"alg\":\"HS256\"}")
	if header == "" {
		return
	}

	mdStr := encoder.Sha256(header + payload + secret)

	return header + "." + payload + "." + mdStr
}

// errCode 1:格式不正确 2:签名错误 3:payload错误 4:过期
const (
	FormatError    = 1 + iota
	SignatureError
	PayloadError
	ExpiredError
)

func ErrcodeString(code int) string {
	return map[int]string{
		0: "success",
		1: "FormatError",
		2: "SignatureError",
		3: "PayloadError",
		4: "ExpiredError",
	}[code]
}

func JWTDecode(secret,in string) (jwtData JWTData, errCode int) {
	data := strings.Split(in, ".")
	jwtData = JWTData{}

	if len(data) != 3 {
		errCode = FormatError
		return
	}
	header := data[0]
	payload := data[1]
	sign := data[2]
	mdStr := encoder.Sha256(header + payload + secret)

	if mdStr != sign {
		errCode = SignatureError
		return
	}

	dataJson := encoder.Base64DecodeString(data[1])

	err := json.Unmarshal([]byte(dataJson), &jwtData)
	if err != nil {
		errCode = PayloadError
		return
	}

	if time.Now().Unix() > jwtData.Exp {
		errCode = ExpiredError
		return
	}

	return
}
