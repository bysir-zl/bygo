package wx_mch

import (
	"testing"
	"github.com/bysir-zl/bygo/util/payment/core"
)

func TestSign(t *testing.T) {
	req := transfersReq{
		NonceStr:       core.RandStr(),
		PartnerTradeNo:"1",
		Openid:         "openid",
		CheckName:      "cname",
		ReUserName:     "uname",
		Amount:         1,
		Desc:           "desc",
		SpbillCreateIp: "1",

		MchAppid: "123",
		Mchid:    "123",
	}

	// 签名
	req.Sign = SignData(req, mckKey)
	t.Log(req.Sign)
}
