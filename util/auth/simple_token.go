package auth

import (
	"errors"
	"fmt"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/encoder"
	"strconv"
	"strings"
	"time"
)

// token生成机制：使用uid+时间+随机数 并MD5 做为Value，将用户组+用户id 做为Key ：存入Redis。返回token为：用户id@Value
// 验证：从token取出用户id和Value再 使用用户组+用户id 做为Key从redis读取到ValueRd，如果Value==ValueRd则成功

// 用法
/*
生成
token,value := CreateToken
key := CreateKey
cache.Set(key,value)
验证
key := GetKeyByToken
cacheToken := cache.Get(key)
VerifyToken(cacheToken,token)
 */

// 生成一个token
func CreateToken(uid int, types string) (token string, value string) {
	key := fmt.Sprintf("%s@%d@%d@%d", types, uid, time.Now().UnixNano(), util.Rand(0, 999999))
	value = encoder.Md5String(key)
	token = fmt.Sprintf("%d@%s", uid, value)
	return
}

// 生成一个key
func CreateKey(uid int, types string) (cacheKey string) {
	cacheKey = fmt.Sprintf("%s%d", types, uid)
	return
}

// 验证时通过token获得缓存key
func GetKeyByToken(token, types string) (cacheKey string, err error) {
	uidAndToken := strings.Split(token, "@")
	if len(uidAndToken) == 1 || uidAndToken[1] == "" {
		err = errors.New("token format error")
		return
	}
	uidStr := uidAndToken[0]
	cacheKey = fmt.Sprintf("%s%s", types, uidStr)
	return
}

// 验证缓存Token与用户Token是否匹配
func VerifyToken(cacheToken string, userToken string) (uid int, err error) {
	uidAndToken := strings.Split(userToken, "@")
	if len(uidAndToken) == 1 || uidAndToken[1] == "" {
		err = errors.New("token format error")
		return
	}
	uidStr := uidAndToken[0]
	verify := uidAndToken[1]

	if err != nil {
		return
	}
	if verify != cacheToken {
		err = errors.New("token invalid")
		return
	}
	uid, err = strconv.Atoi(uidStr)
	return
}
