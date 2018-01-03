package wx_mch

var (
	mchAppid    string // 微信分配的账号ID（企业号corpid即为此appId）
	mchid       string // 微信支付分配的商户号
	mckKey      string
	certPemPath = "./wechat/cert.pem"
	keyPemPath  = "./wechat/key.pem"
	cAPemPath   = "./wechat/ca.pem"
)

// 请先初始化
// 几个path请使用绝对path 或者 运行二进制的相对path
func Init(MchAppid, Mchid, MckKey, CertPemPath, KeyPemPath, CAPemPath string) {
	mchAppid = MchAppid
	mchid = Mchid
	mckKey = MckKey
	certPemPath = CertPemPath
	keyPemPath = KeyPemPath
	cAPemPath = CAPemPath
}
