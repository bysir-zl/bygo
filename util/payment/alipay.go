package payment

import (
	"strconv"
	"github.com/bysir-zl/bygo/util/payment/core"
)

var (
	aliPayNotifyUrl = ""
)

type AliPayNotify struct {
	TradeNO, OutTradeNO string
	Amount              float64
}

// 支付宝调起支付准备
func CreateAliPayPayInfo(tradeNo, subject, totalFee, body string) string {
	q := core.NewAPPayReqForApp()
	q.OutTradeNO = tradeNo
	q.Subject = subject
	q.TotalFee = totalFee
	q.Body = body

	q.NotifyURL = aliPayNotifyUrl

	return q.String()
}

func CheckAliPayNotify(request []byte) (data AliPayNotify, err error) {
	q, err := core.NewAPPayResultNotifyArgs(request)
	if err != nil {
		return
	}
	amount, _ := strconv.ParseFloat(q.TotalFee, 64)

	data = AliPayNotify{
		TradeNO:    q.TradeNO,
		Amount:     amount,
		OutTradeNO: q.OutTradeNO,
	}

	return
}

func InitAli(alipayKey, alipayPartner, alipaySellerEmail, privateKey, publicKey, notifyUrl string) {
	aliPayNotifyUrl = notifyUrl

	config := core.APKeyConfig{
		ALIPAY_KEY:          alipayKey,
		PARTNER_ID:          alipayPartner,
		SELLER_EMAIL:        alipaySellerEmail,
		PARTNET_PRIVATE_KEY: privateKey,
		ALIPAY_PUBLIC_KEY:   publicKey,
	}
	core.InitAPKey(config)
}
