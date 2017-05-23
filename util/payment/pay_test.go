package payment

import "testing"

func TestBbnpay(t *testing.T) {
	c := BbnPayConfig{
		Key:   "",
		AppId: "",
	}
	p := NewBbnPay(c)
	i := BbnPayPlaceOrder{
		Money:     1,
		GoodsId:   153,
		NotifyUrl: "123",
		PcorderId: "order_1",
		PcuserId:  "1",
		GoodsName: "test",
	}
	payInfo, err := p.PlaceOrder(&i)
	t.Log(err, payInfo)
}

func TestAlipay(t *testing.T) {
	config := APKeyConfig{
		ALIPAY_KEY:   "key",
		PARTNER_ID:   "alipay_partner",
		SELLER_EMAIL: "alipay_seller_email",
		PARTNET_PRIVATE_KEY: `-----BEGIN RSA PRIVATE KEY-----
		---
-----END RSA PRIVATE KEY-----
`,
		ALIPAY_PUBLIC_KEY: `-----BEGIN PUBLIC KEY-----
		---
-----END PUBLIC KEY-----`,
	}
	InitAPKey(config)
	q := NewAPPayReqForApp()
	q.OutTradeNO = "1"
	q.Subject = "tital"
	q.TotalFee = "1"
	q.Body = "body"

	q.NotifyURL = "123"
	t.Log(q.String())
}
