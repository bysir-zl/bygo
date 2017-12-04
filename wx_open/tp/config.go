package tp

var (
	Token     string
	AesKey    string
	AppId     string
	AppSecret string
)

// 在调用第三方平台的时候请务必Init
func InitThirdParty(token, aeskey, appid, appSecret string) {
	Token = token
	AesKey = aeskey
	AppId = appid
	AppSecret = appSecret
}
