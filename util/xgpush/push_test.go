package xgpush

import (
	"github.com/bysir-zl/bygo/log"
	"testing"
)

const testToken = "d4c181ebdc1bb8ea4d58ae72352c29fffcab011c"
const (
	access_id_android  = "2100250859"
	access_key_android = "AZ5ERH53Y75N"
	secret_key_android = "349dee0c5dc896d2467ef65a0debcc2d"

	access_id_ios  = "2100250860"
	access_key_ios = "I9U4C34B2KHL"
	secret_key_ios = "ccb870b95168d772d234c2eee19b70c6"
)

func TestAndroid_PushSingle(t *testing.T) {
	per := NewPusherAndroid(access_id_android, access_key_android, secret_key_android)
	mess := Message{
		MessageForAndroid: MessageForAndroid{
			Title:   "title",
			Content: "content",
		},
	}
	err := per.PushSingleDriver(testToken, 1, "", mess)
	if err != nil {
		t.Error(err)
	}
}

// 设置标签, 但是好像不生效
func TestPushAndroid_SetTags4Token(t *testing.T) {
	per := NewPusherAndroid(access_id_android, access_key_android, secret_key_android)
	err := per.SetTags4Token([][2]string{{"G1", testToken}})
	if err != nil {
		t.Error(err)
	}
}

// 获取标签,目前是获取不到的
func TestPushAndroid_GetTokenTags(t *testing.T) {
	per := NewPusherAndroid(access_id_android, access_key_android, secret_key_android)
	r, err := per.GetTokenTags(testToken)
	if err != nil {
		t.Error(err)
	}
	log.Info("test", r)
}

// 所以这个也没法测试
func TestPushAndroid_PushByTags(t *testing.T) {
	per := NewPusherAndroid(access_id_android, access_key_android, secret_key_android)
	mess := Message{
		MessageForAndroid: MessageForAndroid{
			Title:   "title",
			Content: "content",
		},
	}
	err := per.PushByTags([]string{"G1"}, "OR", 1, "", mess)
	if err != nil {
		t.Error(err)
	}
}
