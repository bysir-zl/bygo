package xgpush

import (
	"encoding/json"
)

// 腾讯push

const (
	host    = "openapi.xg.qq.com"
	version = "v2"
)

type AtyAttr struct {
	If int `json:"if,omitempty"` // 创 建 通 知 时 ， intent 的 属 性
	Pf int `json:"pf,omitempty"` // PendingIntent 的属性，
}

type Browser struct {
	Url     string `json:"url,omitempty"` // url
	Confirm int `json:"confirm"`          // 是否需要用户确认
}

type Action struct {
	ActionType int `json:"action_type,omitempty"` // 动作类型，1 打开 activity 或 app 本身，2 打开浏览器，3 打开 Intent
	Activity   string `json:"activity,omitempty"`
	AtyAttr    *AtyAttr `json:"aty_attr,omitempty"` // activity 属性，只针对 action_type=1 的情况
	Browser    *Browser `json:"browser,omitempty"`  //
	Intent     string `json:"intent,omitempty"`     // 打开 Intent
}

type MessageForAndroid struct {
	Content       string `json:"content"`
	Title         string `json:"title"`
	NId           int `json:"n_id,omitempty"`                         // 通知 id，选填。若大于 0，则会覆盖先前弹出的相同 id 通知；若为 0，展示 本条通知且丌影响其他通知；若为-1，将清除先前弹出的所有通知，仅展示本条通知。默认为 0
	BuilderId     int `json:"builder_id"`                             // 本地通知样式，必填
	Ring          int `json:"ring"`                                   // 是否响铃，0 否，1 是，下同。选填，默认 1
	RingRaw       string `json:"ring_raw,omitempty"`                  // 指定应用内的声音（ring.mp3），选填
	Vibrate       int `json:"vibrate"`                                // 是否振动，选填，默认 1
	Lights        int `json:"lights"`                                 // 是否呼吸灯，
	Clearable     int `json:"clearable"`                              // 通知栏是否可清除，选填，默认 1
	IconType      int `json:"icon_type,omitempty"`                    // 默认 0，通知栏图标是应用内图标还是上传图标,0 是应用内图标，1 是上 传图标,选填
	IconRes       string `json:"icon_res,omitempty"`                  // 应用内图标文件名（xg.png）或者下载图标的 url 地址，选填
	StyleId       int `json:"style_id,omitempty"`                     // 应用内图标文件名（xg.png）或者下载图标的 url 地址，选填
	SmallIcon     string `json:"small_icon,omitempty"`                // 指定状态栏的小图片(xg.png),选填
	Action        *Action `json:"action,omitempty"`                   // 选填。默认为打开 app
	CustomContent map[string]string `json:"custom_content,omitempty"` // 用户自定义的 key-value，选填
}

type MessageForIos struct {
	Aps map[string]interface{} `json:"aps"`
}

type Message struct {
	MessageForAndroid
	MessageForIos
}

func (p Message) StringAndroid() string {
	m, _ := json.Marshal(p.MessageForAndroid)
	return string(m)
}

func (p Message) StringIos() string {
	m, _ := json.Marshal(p.MessageForIos)
	return string(m)
}

func (p Message) IsEmptyAndroid() bool {
	return p.MessageForAndroid.Title == ""
}

func (p Message) IsEmptyIos() bool {
	return p.MessageForIos.Aps == nil || len(p.MessageForIos.Aps) == 0
}

type Result struct {
	RetCode int `json:"ret_code,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

type Pusher interface {
	PushSingleDriver(token string, typ int, sendTime string, message Message) (err error)
	PushSingleAccount(account string, typ int, sendTime string, message Message) (err error)
	PushAll(typ int, sendTime string, message Message) (err error)
	PushByTags(tags []string, tagsOpType string, typ int, sendTime string, message Message) (err error)
	SetTags4Token(tags [][2]string) (err error)
	DelTags4Token(tags [][2]string) (err error)
	GetTokenTags(token string) (err error, tags []string)
}
