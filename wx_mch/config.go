package wx_mch

var (
	MchAppid string // 微信分配的账号ID（企业号corpid即为此appId）
	Mchid string // 微信支付分配的商户号
	MckKey string
	CertPemPath = "./wechat/cert.pem"
	KeyPemPath = "./wechat/key.pem"
	CAPemPath = "./wechat/ca.pem"
)