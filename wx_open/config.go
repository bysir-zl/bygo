package wx_open

var (
	Token     string
	AesKey    string
	AppId     string
	AppSecret string
)

// 请务必Init
func Init(token, aeskey, appid, appSecret string) {
	Token = token
	AesKey = aeskey
	AppId = appid
	AppSecret = appSecret
}
