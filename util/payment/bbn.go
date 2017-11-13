package payment

import (
	"strconv"
	"github.com/bysir-zl/bygo/util/payment/core"
)

var (
	bbnPayNotifyUrl = ""

	bbnPay        *core.BbnPay
	bbnPayGoodsId = 0
)

type BbnPayNotifyClient struct {
	TradeNO, OutTradeNO string
	Amount              float64
}

// bbn调起支付准备
func CreateBbnPayInfo(GoodsName, PcorderId, PcuserId string, money string) (payInfo string, err error) {
	m, _ := strconv.ParseFloat(money, 64)
	mInt := int(m * 100)

	i := core.BbnPayPlaceOrder{
		Money:     mInt,
		GoodsId:   bbnPayGoodsId,
		GoodsName: GoodsName,
		PcorderId: PcorderId,
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

	myTradeNo := q.Cporderid

	response.TradeNO = q.Transid
	response.OutTradeNO = myTradeNo
	response.Amount = float64(q.Money) / 100
	return
}

func InitBbn(key, appid string, goodsId int, notifyUrl string) {
	bbnPayNotifyUrl = notifyUrl
	bbnPayGoodsId = goodsId
	c := core.BbnPayConfig{
		Key:   key,
		AppId: appid,
	}
	bbnPay = core.NewBbnPay(c)
}
