package simple

// 简易支付客户端

import (
	"errors"
	"github.com/bysir-zl/bygo/util/payment"
	"strconv"
	"strings"
)

var (
	aliPayNotifyUrl = ""
	bbnPayNotifyUrl = ""

	bbnPay        *payment.BbnPay
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

var typeSp = "@.@"

// 支付宝调起支付准备
func CreateAliPayPayInfo(tradeNo, subject, totalFee, body, types string) string {
	q := payment.NewAPPayReqForApp()
	q.OutTradeNO = types + typeSp + tradeNo
	q.Subject = subject
	q.TotalFee = totalFee
	q.Body = body

	q.NotifyURL = aliPayNotifyUrl

	return q.String()
}

func CheckAliPayNotify(request []byte) (data AliPayNotifyClient, err error) {
	q, err := payment.NewAPPayResultNotifyArgs(request)
	if err != nil {
		return
	}
	amount, _ := strconv.ParseFloat(q.TotalFee, 64)
	ps := strings.Split(q.OutTradeNO, typeSp)
	if len(ps) != 2 {
		err = errors.New("outTradeNo format error: " + q.OutTradeNO)
		return
	}
	types := ps[0]
	myTradeNo := ps[1]

	data = AliPayNotifyClient{
		TradeNO:    q.TradeNO,
		Amount:     amount,
		OutTradeNO: myTradeNo,
		Type:       types,
	}

	return
}

// 微信调起支付准备
func CreateBbnPayInfo(GoodsName, PcorderId, PcuserId string, money, types string) (payInfo string, err error) {
	m, _ := strconv.ParseFloat(money, 64)
	mInt := int(m * 100)

	i := payment.BbnPayPlaceOrder{
		Money:     mInt,
		GoodsId:   bbnPayGoodsId,
		GoodsName: GoodsName,
		PcorderId: types + typeSp + PcorderId,
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
	q, err := bbnPay.Notify(data, sign)
	if err != nil {
		return
	}

	ps := strings.Split(q.Cporderid, typeSp)
	if len(ps) != 2 {
		err = errors.New("outTradeNo format error: " + q.Cporderid)
		return
	}
	types := ps[0]
	myTradeNo := ps[1]

	response.TradeNO = q.Transid
	response.OutTradeNO = myTradeNo
	response.Type = types
	response.Amount = float64(q.Money) / 100
	return
}

func InitBbn(key, appid string, goodsId int, notifyUrl string) {
	bbnPayNotifyUrl = notifyUrl
	bbnPayGoodsId = goodsId
	c := payment.BbnPayConfig{
		Key:   key,
		AppId: appid,
	}
	bbnPay = payment.NewBbnPay(c)
}

func InitAli(alipay_key, alipay_partner, alipay_seller_email, privateKey, publicKey, notifyUrl string) {
	aliPayNotifyUrl = notifyUrl

	config := payment.APKeyConfig{
		ALIPAY_KEY:          alipay_key,
		PARTNER_ID:          alipay_partner,
		SELLER_EMAIL:        alipay_seller_email,
		PARTNET_PRIVATE_KEY: privateKey,
		ALIPAY_PUBLIC_KEY:   publicKey,
	}
	payment.InitAPKey(config)
}
