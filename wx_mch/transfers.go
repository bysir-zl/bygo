package wx_mch

import (
	"encoding/xml"
	"errors"
)

// 企业付款参数
// https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
type TransfersParams struct {
	NonceStr       string // 随机字符串，不长于32位
	PartnerTradeNo string // 商户订单号，需保持唯一性 (只能是字母或者数字，不能包含有符号)
	Openid         string // 商户appid下，某用户的openid
	CheckName      string // NO_CHECK：不校验真实姓名 FORCE_CHECK：强校验真实姓名
	ReUserName     string // 可选, 收款用户真实姓名。如果check_name设置为FORCE_CHECK，则必填用户真实姓名
	Amount         int    // 企业付款金额，单位为分
	Desc           string // 企业付款操作说明信息。必填。
	SpbillCreateIp string // 调用接口的机器Ip地址
}

const (
	TransfersParamsCheckNameNoCheck    = "NO_CHECK"
	TransfersParamsCheckNameForceCheck = "FORCE_CHECK"
)

type (
	transfersReq struct {
		XMLName        struct{} `xml:"xml" sign:"false"`              // root node name
		MchAppid       string   `xml:"mch_appid" sign:"true"`        // 微信分配的账号ID（企业号corpid即为此appId）
		Mchid          string   `xml:"mchid" sign:"true"`            // 微信支付分配的商户号
		NonceStr       string   `xml:"nonce_str" sign:"true"`        // 随机字符串，不长于32位
		PartnerTradeNo string   `xml:"partner_trade_no" sign:"true"` // 商户订单号，需保持唯一性 (只能是字母或者数字，不能包含有符号)
		Openid         string   `xml:"openid" sign:"true"`           // 商户appid下，某用户的openid
		CheckName      string   `xml:"check_name" sign:"true"`       // NO_CHECK：不校验真实姓名 FORCE_CHECK：强校验真实姓名
		ReUserName     string   `xml:"re_user_name" sign:"true"`     // 可选, 收款用户真实姓名。如果check_name设置为FORCE_CHECK，则必填用户真实姓名
		Amount         int      `xml:"amount" sign:"true"`           // 企业付款金额，单位为分
		Desc           string   `xml:"desc" sign:"true"`             // 企业付款操作说明信息。必填。
		SpbillCreateIp string   `xml:"spbill_create_ip" sign:"true"` // 调用接口的机器Ip地址
		Sign           string   `xml:"sign" sign:"false"`            // 签名
	}
	TransfersRsp struct {
		XMLName        struct{} `xml:"xml"`
		ReturnCode     string   `xml:"return_code"`
		ReturnMsg      string   `xml:"return_msg"`
		MchAppid       string   `xml:"mch_appid"`
		Mchid          string   `xml:"mchid"`
		DeviceInfo     string   `xml:"device_info"`
		NonceStr       string   `xml:"nonce_str"`
		ResultCode     string   `xml:"result_code"`
		PartnerTradeNo string   `xml:"partner_trade_no"`
		PaymentNo      string   `xml:"payment_no"`
		PaymentTime    string   `xml:"payment_time"`
	}
)

func Transfers(p TransfersParams) (rsp TransfersRsp, err error) {
	req := transfersReq{
		NonceStr:       p.NonceStr,
		PartnerTradeNo: p.PartnerTradeNo,
		Openid:         p.Openid,
		CheckName:      p.CheckName,
		ReUserName:     p.ReUserName,
		Amount:         p.Amount,
		Desc:           p.Desc,
		SpbillCreateIp: p.SpbillCreateIp,

		MchAppid: mchAppid,
		Mchid:    mchid,
	}

	// 签名
	req.Sign = SignData(req, mckKey)
	post, _ := xml.Marshal(req)
	data, err := SecurePost("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", post)
	if err != nil {
		return
	}

	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return rsp, errors.New("xml Unmarshal err:" + err.Error())
	}

	return rsp, nil
}
