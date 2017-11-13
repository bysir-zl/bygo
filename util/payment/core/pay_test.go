package core

import (
	"testing"
	"encoding/xml"
)

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

func TestWxPay(t *testing.T) {
	wxConfig := WXKeyConfig{
		APP_ID:  "123",
		MCH_ID:  "123",
		MCH_KEY: "123",
	}

	InitWXKey(wxConfig)

	o := WXUnifiedOrderRequest{
		Body:           "hello",
		OutTradeNo:     "orderId1",
		TotalFee:       "1", // 1分
		SpBillCreateIp: "123.12.12.123",
		NotifyURL:      "http://123123/",
		TradeType:      "JSAPI",
		OpenId:         "123",
	}

	rsp, err := o.Post()
	if err != nil {
		t.Fatal(err)
	}
	err = rsp.Error()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp.PrePayId)
}

func TestWxPayCallback(t *testing.T) {
	wxConfig := WXKeyConfig{
		APP_ID:  "123",
		MCH_ID:  "123",
		MCH_KEY: "123",
	}

	InitWXKey(wxConfig)

	body := []byte(`
	<xml>
   <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
   <attach><![CDATA[支付测试]]></attach>
   <bank_type><![CDATA[CFT]]></bank_type>
   <fee_type><![CDATA[CNY]]></fee_type>
   <is_subscribe><![CDATA[Y]]></is_subscribe>
   <mch_id><![CDATA[10000100]]></mch_id>
   <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
   <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
   <out_trade_no><![CDATA[1409811653]]></out_trade_no>
   <result_code><![CDATA[SUCCESS]]></result_code>
   <return_code><![CDATA[SUCCESS]]></return_code>
   <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
   <sub_mch_id><![CDATA[10000100]]></sub_mch_id>
   <time_end><![CDATA[20140903131540]]></time_end>
   <total_fee>1</total_fee>
   <trade_type><![CDATA[JSAPI]]></trade_type>
   <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
 </xml> `)

	n := PayNotify{}
	xml.Unmarshal(body, &n)
	err := n.IsError()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n.OutTradeNo, n.TotalFee)
}
