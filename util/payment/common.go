package payment

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	NET_ERROR            = errors.New("net error")
	DATA_UNMARSHAL_ERROR = errors.New("data unmarshal error")
)

var rs string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//获得当前时间戳
func TimeNow() int64 {
	return time.Now().Unix()
}

//获得当前时间戳
func TimeNowString() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func TimeString(d int64) string {
	return time.Now().Add(time.Duration(d) * time.Second).Format("20060102150405")
}

//获得随机字符串
func RandStr() string {
	s := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 32; i++ {
		v := rand.Int() % len(rs)
		s += string(rs[v])
	}
	return s
}
