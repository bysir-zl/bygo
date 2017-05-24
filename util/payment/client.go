package payment

// 简易支付客户端

import (
	"errors"
	"strconv"
)

var (
	aliPayNotifyUrl = ""
	bbnPayNotifyUrl = ""

	bbnPay        *BbnPay
	bbnPayGoodsId = 0
)

type AliPayNotifyClient struct {
	TradeNO, OutTradeNO string
	Amount              float64

	Type string // 类型: recharge
}

type BbnPayNotifyClient struct {
	TradeNO, OutTradeNO string
	Amount              float64

	Type string // 类型: recharge
}

const (
	Type_Buy = "re_"
)

// 支付宝调起支付准备
func CreateAliPayPayInfo(tradeNo, subject, totalFee, body string) string {
	q := NewAPPayReqForApp()
	q.OutTradeNO = Type_Buy + tradeNo
	q.Subject = subject
	q.TotalFee = totalFee
	q.Body = body

	q.NotifyURL = aliPayNotifyUrl

	return q.String()
}

func CheckAliPayNotify(request []byte) (data AliPayNotifyClient, err error) {
	q, err := NewAPPayResultNotifyArgs(request)
	if err != nil {
		return
	}
	amount, _ := strconv.ParseFloat(q.TotalFee, 64)
	t := ""
	myTradeNo := q.OutTradeNO[3:]
	switch q.OutTradeNO[:3] {
	case Type_Buy:
		t = Type_Buy
	default:
		err = errors.New("error type," + q.OutTradeNO[:3])
		return
	}

	data = AliPayNotifyClient{
		TradeNO:    q.TradeNO,
		Amount:     amount,
		OutTradeNO: myTradeNo,
		Type:       t,
	}

	return
}

// 微信调起支付准备
func CreateBbnPayInfo(GoodsName, PcorderId, PcuserId string, money string) (payInfo string, err error) {
	m, _ := strconv.ParseFloat(money, 64)
	mInt := int(m * 100)

	i := BbnPayPlaceOrder{
		Money:     mInt,
		GoodsId:   bbnPayGoodsId,
		GoodsName: GoodsName,
		PcorderId: Type_Buy + PcorderId,
		NotifyUrl: bbnPayNotifyUrl,
		PcuserId:  PcuserId,
	}
	payInfo, err = bbnPay.PlaceOrder(&i)
	if err != nil {
		return
	}

	return
}

func CheckBbnPayNotify(data, sign string) (response BbnPayNotifyClient, err error) {
	bp, err := bbnPay.Notify(data, sign)
	if err != nil {
		return
	}

	t := ""
	myTradeNo := bp.Cporderid[3:]
	switch bp.Cporderid[:3] {
	case Type_Buy:
		t = Type_Buy
	default:
		err = errors.New("error type," + bp.Cporderid[:3])
		return
	}

	response.TradeNO = bp.Transid
	response.OutTradeNO = myTradeNo
	response.Type = t
	response.Amount = float64(bp.Money) / 100
	return
}

func IninBbn(key, appid, notifyUrl string, goodsId int) {
	bbnPayNotifyUrl = notifyUrl
	bbnPayGoodsId = goodsId
	c := BbnPayConfig{
		Key:   key,
		AppId: appid,
	}
	bbnPay = NewBbnPay(c)
}

func InitAli(alipay_key, alipay_partner, alipay_seller_email, privateKey, publicKey, notifyUrl string) {
	aliPayNotifyUrl = notifyUrl

	config := APKeyConfig{
		ALIPAY_KEY:          alipay_key,
		PARTNER_ID:          alipay_partner,
		SELLER_EMAIL:        alipay_seller_email,
		PARTNET_PRIVATE_KEY: privateKey,
		ALIPAY_PUBLIC_KEY:   publicKey,
	}
	InitAPKey(config)
}
